package common

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Jille/easymutex"
	"github.com/Jille/genericz/mapz"
	pb "github.com/SnoozeThis-org/logwait/proto"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var (
	ObserverAddress = pflag.String("observer-address", "localhost:1600", "gRPC address of the observer")
)

type Observable struct {
	observable *pb.Observable
	Regexps    map[string]*regexp.Regexp
	observed   bool
}

type Service struct {
	StartObserving             func()
	StopObserving              func()
	scannerType                string
	mtx                        sync.Mutex
	observables                map[string]Observable
	observerClient             pb.ObserverServiceClient
	observerStream             pb.ObserverService_CommunicateClient
	receivedInitialObservables chan struct{}
	filterableFields           []string
}

func NewService(c *grpc.ClientConn, scannerType string) *Service {
	return &Service{
		scannerType:                scannerType,
		observerClient:             pb.NewObserverServiceClient(c),
		observables:                map[string]Observable{},
		receivedInitialObservables: make(chan struct{}),
	}
}

func (s *Service) ConnectLoop() {
	failures := 0
	for {
		if err := s.connectToObserver(); err != nil {
			log.Printf("Error while talking to observer: %v", err)
		}
		failures++
		time.Sleep(time.Duration(failures) * time.Second)
	}
}

func (s *Service) MatchObservables(f func(o Observable) bool) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for id, o := range s.observables {
		if o.observed {
			continue
		}
		if f(o) {
			o.observed = true
			s.observables[id] = o
			go s.notifyObserver(id)
		}
	}
}

func (s *Service) connectToObserver() error {
	em := easymutex.EasyMutex{L: &s.mtx}
	defer em.Unlock() // Drop the lock if we happen to return while holding it.
	ctx := context.Background()
	stream, err := s.observerClient.Communicate(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_Register{
			Register: &pb.RegisterScannerRequest{
				ScannerType: s.scannerType,
			},
		},
	}); err != nil {
		return err
	}
	msg, err := stream.Recv()
	if err != nil {
		return err
	}

	em.Lock()
	s.observerStream = stream

	r := msg.GetMsg().(*pb.ObserverToScanner_Register).Register
	for _, o := range r.ActiveObservables {
		s.addObservable_locked(o)
	}
	em.Unlock()

	defer func() {
		em.Lock()
		s.observerStream = nil
	}()

	select {
	case <-s.receivedInitialObservables:
	default:
		close(s.receivedInitialObservables)
	}

	em.Lock()
	if len(s.filterableFields) > 0 {
		if err := s.observerStream.Send(&pb.ScannerToObserver{
			Msg: &pb.ScannerToObserver_UpdateFields{
				UpdateFields: &pb.UpdateFields{
					Fields: s.filterableFields,
				},
			},
		}); err != nil {
			return err
		}
	}
	em.Unlock()

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		switch m := msg.Msg.(type) {
		case *pb.ObserverToScanner_NewObservable:
			s.addObservable(m.NewObservable)
		case *pb.ObserverToScanner_CancelObservable:
			s.cancelObservable(m.CancelObservable)
		case *pb.ObserverToScanner_TestFilters:
			reply := s.testFilters(m.TestFilters)
			reply.TestId = m.TestFilters.TestId
			if err := stream.Send(&pb.ScannerToObserver{
				Msg: &pb.ScannerToObserver_TestFilters{
					TestFilters: reply,
				},
			}); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected %T from Scanner", msg.Msg)
		}
	}
}

func (s *Service) addObservable(o *pb.Observable) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.addObservable_locked(o)
}

func (s *Service) addObservable_locked(o *pb.Observable) {
	val, ok := s.observables[o.Id]

	if ok {
		if val.observed {
			// The Observer just asked us to wait for this observable, but we've already seen it.
			// That probably means our previous signal of that was lost and we should resend it.
			// This is likely to occur during reconnection to the Observer when the Observer sends all active Observables again.
			go s.notifyObserver(o.Id)
		}
		return
	}

	val.observable = o
	val.Regexps = make(map[string]*regexp.Regexp, len(o.Filters))
	for _, f := range o.Filters {
		re, err := regexp.Compile(f.Regexp)
		if err != nil {
			s.sendRejection(o.Id, "invalid regexp for field "+f.Field+": "+err.Error())
			return
		}
		val.Regexps[f.Field] = re
	}

	if s.StartObserving != nil && len(s.observables) == 0 {
		s.StartObserving()
	}
	s.observables[o.Id] = val
}

func (s *Service) cancelObservable(id string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if len(s.observables) == 0 {
		// Avoid calling StopObserving below if we got a cancel for a non-existent id and were already stopped.
		return
	}

	delete(s.observables, id)

	if s.StopObserving != nil && len(s.observables) == 0 {
		s.StopObserving()
	}
}

func (s *Service) notifyObserver(id string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	o, ok := s.observables[id]
	if !ok {
		return
	}
	if !o.observed {
		return
	}
	if s.observerStream == nil {
		// We'll send the message once we get reconnected.
		return
	}
	// We ignore any error from Send(). The Recv() thread will soon notice the connection being gone too and will reconnect and then we'll retry sending this observation.
	_ = s.observerStream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_ObservedObservable{
			ObservedObservable: id,
		},
	})
}

func (s *Service) sendRejection(id, msg string) {
	// We ignore any error from Send(). The Recv() thread will soon notice the connection being gone too and will reconnect and then we'll retry sending this observation.
	_ = s.observerStream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_RejectObservable{
			RejectObservable: &pb.RejectObservable{
				Id:      id,
				Message: msg,
			},
		},
	})
}

func (s *Service) testFilters(req *pb.TestFiltersRequest) *pb.TestFiltersResponse {
	var ret pb.TestFiltersResponse
	for _, f := range req.Filters {
		_, err := regexp.Compile(f.Regexp)
		if err != nil {
			ret.Errors = append(ret.Errors, "invalid regexp for field "+f.Field+": "+err.Error())
		}
	}
	return &ret
}

// WaitForInitialConnection waits the filters to be received from the Observer.
// This ensures that existing Observables don't miss any logs even if the scanner otherwise keeps track of its log position.
func (s *Service) WaitForInitialConnection() {
	<-s.receivedInitialObservables
}

func (s *Service) SetFilterableFields(fields []string) {
	s.mtx.Lock()
	s.filterableFields = fields
	if s.observerStream == nil {
		s.mtx.Unlock()
		// We'll send the message once we get (re)connected.
		return
	}
	s.mtx.Unlock()
	// We ignore any error from Send(). The Recv() thread will soon notice the connection being gone too and will reconnect and then we'll retry sending this observation.
	_ = s.observerStream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_UpdateFields{
			UpdateFields: &pb.UpdateFields{
				Fields: fields,
			},
		},
	})
}

func (s *Service) CreateFieldLearner(separator byte) *FieldLearner {
	return &FieldLearner{
		parent:    s,
		separator: separator,
	}
}

type FieldLearner struct {
	parent    *Service
	seen      mapz.SyncMap[string, struct{}]
	ignored   mapz.SyncMap[string, struct{}]
	separator byte

	addMtx sync.Mutex
}

func (l *FieldLearner) Seen(field string) {
	if _, existed := l.seen.Load(field); existed {
		return
	}
	_, ignored := l.checkIgnored(field)
	if !ignored {
		go l.add(field)
	}
}

func (l *FieldLearner) checkIgnored(field string) (parentsLength int, ignored bool) {
	var offset int
	for {
		i := strings.IndexByte(field[offset:], l.separator)
		if i == -1 {
			return offset, false
		}
		dirs := field[:offset+i]
		if _, ignored := l.ignored.Load(dirs); ignored {
			return offset, true
		}
		offset += i + 1
	}
}

func (l *FieldLearner) add(field string) {
	l.addMtx.Lock()
	defer l.addMtx.Unlock()

	parentsLength, ignored := l.checkIgnored(field)
	if ignored {
		return
	}

	if _, existed := l.seen.LoadOrStore(field, struct{}{}); existed {
		return
	}

	fields := make([]string, 0, 100)
	l.seen.Range(func(key string, _ struct{}) bool {
		fields = append(fields, key)
		return true
	})

	if parentsLength > 0 {
		parent := field[:parentsLength+1] // parent plus the separator
		siblings := 0
		for _, f := range fields {
			if strings.HasPrefix(f, parent) {
				siblings++
			}
		}
		if siblings >= 50 {
			l.ignored.Store(field[:parentsLength], struct{}{})
			keep := make([]string, 0, len(fields)-siblings)
			for _, f := range fields {
				if strings.HasPrefix(f, parent) {
					l.seen.Delete(f)
				} else {
					keep = append(keep, f)
				}
			}
			fields = keep
		}
	}
	l.parent.SetFilterableFields(fields)
}
