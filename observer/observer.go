package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/SnoozeThis-org/logwait/config"
	pb "github.com/SnoozeThis-org/logwait/proto"
	"golang.org/x/exp/maps"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	// These settings can also passed through environment variables, e.g. SIGNING_KEY.
	httpPort       = config.FlagSet.Int("http-port", 8080, "Port to serve HTTP interface on")
	grpcPort       = config.FlagSet.Int("grpc-port", 1600, "Port to serve gRPC for the Scanners on")
	token          = config.FlagSet.String("token", "", "The token you got when registering at https://www.snoozethis.com/logs/")
	signingKey     = config.FlagSet.String("signing-key", "", "Secret to sign your filters with")
	allowUnsigned  = config.FlagSet.Bool("allow-unsigned", false, "Whether a signing key is required")
	snoozethisAddr = config.FlagSet.String("snoozethis-addr", "logs.grpc.snoozethis.com:443", "gRPC address of the SnoozeThisLogService")
)

type service struct {
	mtx                       sync.Mutex
	activeObservables         map[string]*pb.Observable
	connectedScanners         map[pb.ObserverService_CommunicateServer]struct{}
	scannerType               string
	fieldsPerScanner          map[pb.ObserverService_CommunicateServer][]string
	snoozeThisClient          pb.SnoozeThisLogServiceClient
	snoozeThisClientStream    pb.SnoozeThisLogService_CommunicateClient
	unconfirmedObservations   map[string]struct{}
	reportedRejections        map[string]struct{}
	initialObservablesFetched chan struct{}
	registeredUrl             string
	pendingTestFilters        map[int64]chan *pb.TestFiltersResponse
	nextTestId                int64
}

func main() {
	config.Parse()

	if *token == "" {
		fmt.Fprintln(os.Stderr, "The token is required (--token or the env var TOKEN). Get one at https://www.snoozethis.com/logs/")
		os.Exit(2)
	}
	if *signingKey == "" && !*allowUnsigned {
		fmt.Fprintln(os.Stderr, "A signing key is required for optimal security. You can waive your security by using --allow-unsigned")
		os.Exit(2)
	}

	c, err := grpc.Dial(*snoozethisAddr, grpc.WithTransportCredentials(credentials.NewTLS(nil)), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time: time.Minute,
	}))
	if err != nil {
		log.Fatalf("Failed to dial %s: %v", *snoozethisAddr, err)
	}

	srv := &service{
		activeObservables:         map[string]*pb.Observable{},
		connectedScanners:         map[pb.ObserverService_CommunicateServer]struct{}{},
		fieldsPerScanner:          map[pb.ObserverService_CommunicateServer][]string{},
		snoozeThisClient:          pb.NewSnoozeThisLogServiceClient(c),
		unconfirmedObservations:   map[string]struct{}{},
		reportedRejections:        map[string]struct{}{},
		initialObservablesFetched: make(chan struct{}),
		pendingTestFilters:        map[int64]chan *pb.TestFiltersResponse{},
	}

	go srv.talkToSnoozeThisLoop()

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

	srv.registerUI()
	log.Fatalf("HTTP server died: %v", http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}

func (s *service) talkToSnoozeThisLoop() {
	failures := 0
	for {
		if err := s.talkToSnoozeThis(); err != nil {
			log.Printf("Error while talking to SnoozeThis: %v", err)
		}
		failures++
		time.Sleep(time.Duration(failures) * time.Second)
	}
}

func (s *service) talkToSnoozeThis() error {
	ctx := context.Background()
	stream, err := s.snoozeThisClient.Communicate(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&pb.ObserverToSnoozeThis{
		Msg: &pb.ObserverToSnoozeThis_Register{
			Register: &pb.RegisterObserverRequest{
				LogInstanceToken:        *token,
				ObserverProtocolVersion: 20230815,
			},
		},
	}); err != nil {
		return err
	}
	msg, err := stream.Recv()
	if err != nil {
		if status.Code(err) == codes.PermissionDenied {
			log.Fatalf("Your token was rejected")
		}
		return err
	}
	s.mtx.Lock()
	s.snoozeThisClientStream = stream
	s.reportedRejections = map[string]struct{}{}
	r := msg.GetMsg().(*pb.SnoozeThisToObserver_Register).Register
	s.registeredUrl = r.GetRegisteredUrl()
	seen := map[string]struct{}{}
	for _, o := range r.ActiveObservables {
		s.addObservable_locked(o)
		seen[o.GetId()] = struct{}{}
	}
	for id := range s.unconfirmedObservations {
		if _, exists := seen[id]; !exists {
			delete(s.unconfirmedObservations, id)
		}
	}
	seen = nil
	s.mtx.Unlock()

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
			s.mtx.Lock()
			s.addObservable_locked(m.NewObservable)
			s.mtx.Unlock()
		case *pb.SnoozeThisToObserver_CancelObservable:
			s.cancelObservable(m.CancelObservable)
		default:
			// A message we don't know. Most likely a new protocol addition. Let's ignore it.
		}
	}
}

func (s *service) addObservable_locked(o *pb.Observable) {
	if msg := checkObservable(o); msg != "" {
		// If this fails, we'll reconnect soon and SnoozeThis will send us the obstacle again and we'll reject it again.
		_ = s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
			Msg: &pb.ObserverToSnoozeThis_RejectObservable{
				RejectObservable: &pb.RejectObservable{
					Id:      o.Id,
					Message: msg,
				},
			},
		})
		return
	}

	if _, found := s.unconfirmedObservations[o.GetId()]; found {
		_ = s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
			Msg: &pb.ObserverToSnoozeThis_ObservedObservable{
				ObservedObservable: o.GetId(),
			},
		})
		return
	}

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
	delete(s.unconfirmedObservations, id)
	delete(s.reportedRejections, id)
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
	s.unconfirmedObservations[id] = struct{}{}
	if s.snoozeThisClientStream == nil {
		return
	}
	_ = s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
		Msg: &pb.ObserverToSnoozeThis_ObservedObservable{
			ObservedObservable: id,
		},
	})
}

func (s *service) rejectObservable(r *pb.RejectObservable) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if s.snoozeThisClientStream == nil {
		return
	}
	_, reported := s.reportedRejections[r.GetId()]
	if reported {
		// Avoid all scanners sending this message to SnoozeThis.
		return
	}
	if s.snoozeThisClientStream.Send(&pb.ObserverToSnoozeThis{
		Msg: &pb.ObserverToSnoozeThis_RejectObservable{
			RejectObservable: r,
		},
	}) == nil {
		s.reportedRejections[r.GetId()] = struct{}{}
	}
}

func (s *service) Communicate(stream pb.ObserverService_CommunicateServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	r, ok := msg.GetMsg().(*pb.ScannerToObserver_Register)
	if !ok {
		return errors.New("first message from Scanner to observer should be Register")
	}

	<-s.initialObservablesFetched

	s.mtx.Lock()
	if len(s.connectedScanners) == 0 {
		s.scannerType = r.Register.GetScannerType()
	}
	if s.scannerType != r.Register.GetScannerType() {
		others := s.scannerType
		s.mtx.Unlock()
		return fmt.Errorf("All scanners should be the same type. The other scanners are of type %q but you are a %q.", others, r.Register.GetScannerType())
	}
	s.connectedScanners[stream] = struct{}{}
	observables := maps.Values(s.activeObservables)
	s.mtx.Unlock()
	defer func() {
		s.mtx.Lock()
		delete(s.connectedScanners, stream)
		delete(s.fieldsPerScanner, stream)
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
		case *pb.ScannerToObserver_RejectObservable:
			s.rejectObservable(m.RejectObservable)
		case *pb.ScannerToObserver_TestFilters:
			s.handleTestFiltersResponse(m.TestFilters)
		case *pb.ScannerToObserver_UpdateFields:
			s.mtx.Lock()
			s.fieldsPerScanner[stream] = m.UpdateFields.Fields
			s.mtx.Unlock()
		default:
			return fmt.Errorf("unexpected %T from Scanner", msg.Msg)
		}
	}
}

func checkObservable(o *pb.Observable) string {
	if !checkSignature(o) {
		return "signature is invalid"
	}
	return ""
}

func checkSignature(o *pb.Observable) bool {
	if *allowUnsigned {
		return true
	}
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
