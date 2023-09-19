package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/SnoozeThis-org/logwait/config"
	"github.com/SnoozeThis-org/logwait/scanners/common"
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

	srv = common.NewService(c, "http-scanner")
	go srv.ConnectLoop()
	learner = srv.CreateFieldLearner('.')
	srv.WaitForInitialConnection()

	http.ListenAndServe(fmt.Sprintf(":%d", *port), http.HandlerFunc(pathRouter))
}

func pathRouter(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/elasticsearch":
		elasticsearch(w, r)
	case "/elasticsearch/_bulk":
		elasticsearchBulk(w, r)
	case "/ndjson":
		ndjson(w, r)
	default:
		http.NotFound(w, r)
	}
}

func elasticsearch(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		http.Error(w, "Failed to parse JSON body: "+err.Error(), 400)
		return
	}
	handleMessage(m)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte{'{', '}'})
}

type elasticBulkCommand struct {
	Index  *struct{} `json:"index,omitempty"`
	Delete *struct{} `json:"delete,omitempty"`
	Create *struct{} `json:"create,omitempty"`
	Update *struct{} `json:"update,omitempty"`
}

func elasticsearchBulk(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	var ebc elasticBulkCommand
	for {
		if err := d.Decode(&ebc); err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, "Failed to parse JSON body: "+err.Error(), 400)
			return
		}
		if ebc.Delete != nil || ebc.Update != nil {
			http.Error(w, "You can't delete or update logs", 403)
			return
		}
		var m map[string]any
		if err := d.Decode(&m); err != nil {
			http.Error(w, "Failed to parse JSON body: "+err.Error(), 400)
			return
		}
		handleMessage(m)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte{'{', '}'})
}

func ndjson(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	for {
		var m map[string]any
		if err := d.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			http.Error(w, "Failed to parse JSON body: "+err.Error(), 400)
			return
		}
		handleMessage(m)
	}
	w.WriteHeader(204)
	return
}

func handleMessage(m map[string]any) {
	learnAllFields(learner, m, "")
	srv.MatchObservables(func(o common.Observable) bool {
		for field, regexp := range o.Regexps {
			val, ok := jsonPathAsString(m, field)
			if !ok {
				// Value not found or of a non-scalar type.
				return false
			}
			if !regexp.MatchString(val) {
				return false
			}
		}
		return true
	})
}

func jsonPathAsString(m map[string]any, path string) (string, bool) {
	sp := strings.Split(path, ".")
	last := sp[len(sp)-1]
	sp = sp[:len(sp)-1]
	for _, p := range sp {
		sub, ok := m[p].(map[string]any)
		if !ok {
			return "", false
		}
		m = sub
	}
	v, ok := m[last]
	if !ok {
		return "", false
	}
	switch v := v.(type) {
	case string:
		return v, true
	case bool, float64:
		return fmt.Sprint(v), true
	default:
		return "", false
	}
}

func learnAllFields(l *common.FieldLearner, m map[string]any, prefix string) {
	for k, v := range m {
		if prefix != "" {
			k = prefix + "." + k
		}
		sub, ok := v.(map[string]any)
		if ok {
			learnAllFields(l, sub, k)
		} else {
			l.Seen(k)
		}
	}
}
