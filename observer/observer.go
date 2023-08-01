package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"

	pb "github.com/SnoozeThis-org/logwait/proto"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	httpPort   = flag.Int("http-port", 8080, "Port to serve HTTP interface on")
	grpcPort   = flag.Int("grpc-port", 1600, "Port to serve gRPC for the Scanners on")
	signingKey = flag.String("signing-key", "secret", "Secret to sign your filters with")
)

type service struct {
	mtx                       sync.Mutex
	activeObservables         map[string]*pb.Observable
	connectedScanners         map[pb.ObserverService_CommunicateServer]struct{}
	snoozeThisClient          pb.SnoozeThisLogServiceClient
	snoozeThisClientStream    pb.SnoozeThisLogService_CommunicateClient
	initialObservablesFetched chan struct{}
}

func main() {
	flag.Parse()

	c, err := grpc.Dial("https://logs.grpc.snoozethis.io")
	if err != nil {
		log.Fatalf("Failed to dial logs.grpc.snoozethis.io: %v", err)
	}

	srv := &service{
		activeObservables:         map[string]*pb.Observable{},
		connectedScanners:         map[pb.ObserverService_CommunicateServer]struct{}{},
		snoozeThisClient:          pb.NewSnoozeThisLogServiceClient(c),
		initialObservablesFetched: make(chan struct{}),
	}

	go srv.connectLoop()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %d: %v", *grpcPort, err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterObserverServiceServer(s, srv)
	go func() {
		if err := s.Serve(l); err != nil {
			log.Fatalf("gRPC server died: %v", err)
		}
	}()

	log.Fatalf("HTTP server died: %v", http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func (s *service) connectLoop() {
	failures := 0
	for {
		if err := s.connectToSnoozeThis(); err != nil {
			log.Printf("Error while talking to SnoozeThis: %v", err)
		}
		failures++
		time.Sleep(time.Duration(failures) * time.Second)
	}
}

func (s *service) connectToSnoozeThis() error {
	ctx := context.Background()
	stream, err := s.snoozeThisClient.Communicate(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.ObserverToSnoozeThis{
		Msg: &pb.ObserverToSnoozeThis_Register{
			Register: &pb.RegisterObserverRequest{
				LogInstanceToken: "secret",
			},
		},
	}); err != nil {
		return err
	}
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	r := msg.GetMsg().(*pb.SnoozeThisToObserver_Register).Register
	for _, o := range r.ActiveObservables {
		s.addObservable(o)
	}

	select {
	case <-s.initialObservablesFetched:
	default:
		close(s.initialObservablesFetched)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		switch m := msg.Msg.(type) {
		case *pb.SnoozeThisToObserver_NewObservable:
			s.addObservable(m.NewObservable)
		case *pb.SnoozeThisToObserver_CancelObservable:
			s.cancelObservable(m.CancelObservable)
		default:
			return fmt.Errorf("unexpected %T from Scanner", msg.Msg)
		}
	}
}

func (s *service) addObservable(o *pb.Observable) {
	if msg := checkObservable(o); msg != "" {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		if err := s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
			Msg: &pb.ObserverToSnoozeThis_RejectObservable{
				RejectObservable: &pb.RejectObservable{
					Id:      o.Id,
					Message: msg,
				},
			},
		}); err != nil {
			panic(err) // XXX
		}
		return
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.activeObservables[o.GetId()] = o
	for scanner := range s.connectedScanners {
		// If this fails, the Recv() will likely fail soon too and will disconnect this scanner.
		_ = scanner.Send(&pb.ObserverToScanner{
			Msg: &pb.ObserverToScanner_NewObservable{
				NewObservable: o,
			},
		})
	}
}

func (s *service) cancelObservable(id string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.activeObservables, id)
	for scanner := range s.connectedScanners {
		// If this fails, the Recv() will likely fail soon too and will disconnect this scanner.
		_ = scanner.Send(&pb.ObserverToScanner{
			Msg: &pb.ObserverToScanner_CancelObservable{
				CancelObservable: id,
			},
		})
	}
}

func (s *service) observedObservable(id string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if err := s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
		Msg: &pb.ObserverToSnoozeThis_ObservedObservable{
			ObservedObservable: id,
		},
	}); err != nil {
		panic(err) // XXX
	}
}

func (s *service) Communicate(stream pb.ObserverService_CommunicateServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	_, ok := msg.GetMsg().(*pb.ScannerToObserver_Register)
	if !ok {
		return errors.New("first message from Scanner to observer should be Register")
	}

	<-s.initialObservablesFetched

	s.mtx.Lock()
	s.connectedScanners[stream] = struct{}{}
	observables := maps.Values(s.activeObservables)
	s.mtx.Unlock()
	defer func() {
		s.mtx.Lock()
		delete(s.connectedScanners, stream)
		s.mtx.Unlock()
	}()

	if err := stream.Send(&pb.ObserverToScanner{
		Msg: &pb.ObserverToScanner_Register{
			Register: &pb.RegisterScannerResponse{
				ActiveObservables: observables,
			},
		},
	}); err != nil {
		return err
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			return err
		}
		switch m := msg.Msg.(type) {
		case *pb.ScannerToObserver_ObservedObservable:
			s.observedObservable(m.ObservedObservable)
		default:
			return fmt.Errorf("unexpected %T from Scanner", msg.Msg)
		}
	}
}

func checkObservable(o *pb.Observable) string {
	if !checkSignature(o) {
		return "signature is invalid"
	}
	for _, f := range o.Filters {
		if _, err := regexp.Compile(f.GetRegexp()); err != nil {
			return "regexp was invalid"
		}
	}
	return ""
}

func checkSignature(o *pb.Observable) bool {
	expected := calculateSignature(o)
	return hmac.Equal([]byte(expected), []byte(o.GetSignature()))
}

func calculateSignature(o *pb.Observable) string {
	filters := make([]string, len(o.Filters))
	for i, f := range o.Filters {
		filters[i] = fmt.Sprintf("%d%d%s%s", len(f.GetField()), len(f.GetRegexp()), f.GetField(), f.GetRegexp())
	}
	sort.Strings(filters)
	m := hmac.New(sha256.New, []byte(*signingKey))
	for _, f := range filters {
		io.WriteString(m, f)
	}
	sum := m.Sum(nil)
	return hex.EncodeToString(sum[:])
}
