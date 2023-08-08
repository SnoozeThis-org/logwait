package common

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	pb "github.com/SnoozeThis-org/logwait/proto"
	"google.golang.org/grpc"
)

var (
	ObserverAddress = flag.String("observer-address", "localhost:1600", "gRPC address of the observer")
)

type Observable struct {
	observable *pb.Observable
	Regexps    map[string]*regexp.Regexp
	observed   bool
}

type Service struct {
	StartObserving func()
	StopObserving  func()
	mtx            sync.Mutex
	observables    map[string]Observable
	observerClient pb.ObserverServiceClient
	observerStream pb.ObserverService_CommunicateClient
}

func NewService(c *grpc.ClientConn) *Service {
	return &Service{
		observerClient: pb.NewObserverServiceClient(c),
		observables:    map[string]Observable{},
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
	ctx := context.Background()
	stream, err := s.observerClient.Communicate(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_Register{
			Register: &pb.RegisterScannerRequest{},
		},
	}); err != nil {
		return err
	}
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	s.observerStream = stream

	r := msg.GetMsg().(*pb.ObserverToScanner_Register).Register
	for _, o := range r.ActiveObservables {
		s.addObservable(o)
	}

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
		default:
			return fmt.Errorf("unexpected %T from Scanner", msg.Msg)
		}
	}
}

func (s *Service) addObservable(o *pb.Observable) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	val, ok := s.observables[o.Id]

	if ok {
		if val.observed {
			go s.notifyObserver(o.Id)
		}
		return
	}

	val.observable = o
	val.Regexps = make(map[string]*regexp.Regexp, len(o.Filters))
	for _, f := range o.Filters {
		val.Regexps[f.Field] = regexp.MustCompile(f.Regexp)
	}

	if s.StartObserving != nil && len(s.observables) == 0 {
		s.StartObserving()
	}
	s.observables[o.Id] = val
}

func (s *Service) cancelObservable(id string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
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
	s.observerStream.Send(&pb.ScannerToObserver{
		Msg: &pb.ScannerToObserver_ObservedObservable{
			ObservedObservable: id,
		},
	})
}
