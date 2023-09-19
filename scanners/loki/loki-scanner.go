package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/SnoozeThis-org/logwait/config"
	"github.com/SnoozeThis-org/logwait/scanners/common"
	gogo "github.com/gogo/protobuf/proto"
	"github.com/grafana/loki/pkg/push"
	"github.com/klauspost/compress/s2"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var (
	port = pflag.Int("port", 0, "Port to open for HTTP")

	srv     *common.Service
	learner *common.FieldLearner
)

func main() {
	config.Parse()

	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}

	srv = common.NewService(c, "loki-scanner")
	go srv.ConnectLoop()
	learner = srv.CreateFieldLearner('.')
	learner.Seen("line")
	srv.WaitForInitialConnection()

	http.ListenAndServe(fmt.Sprintf(":%d", *port), http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/loki/api/v1/push" {
		http.NotFound(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	body, err = s2.Decode(nil, body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	var msg push.PushRequest
	if err := gogo.Unmarshal(body, &msg); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	for _, s := range msg.Streams {
		for _, e := range s.Entries {
			handleMessage(e.Line, e.StructuredMetadata)
		}
	}

	w.WriteHeader(204)
}

func handleMessage(line string, labels push.LabelsAdapter) {
	for _, l := range labels {
		learner.Seen(l.Name)
	}
	srv.MatchObservables(func(o common.Observable) bool {
	loop:
		for field, regexp := range o.Regexps {
			switch field {
			case "line":
				if !regexp.MatchString(line) {
					return false
				}
			default:
				for _, l := range labels {
					if l.Name == field {
						if !regexp.MatchString(l.Value) {
							return false
						}
						continue loop
					}
				}
				return false
			}
		}
		return true
	})
}
