package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SnoozeThis-org/logwait/scanners/common"
	"google.golang.org/grpc"
)

var (
	readJson = flag.Bool("json", false, "Read json instead of lines")

	srv *common.Service
)

func main() {
	flag.Parse()

	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}
	defer c.Close()

	srv = common.NewService(c)
	if *readJson {
		go jsonConsumer()
	} else {
		go lineConsumer()
	}
	srv.ConnectLoop()
}

func lineConsumer() {
	s := bufio.NewScanner(os.Stdin)
	srv.WaitForInitialConnection()
	for s.Scan() {
		line := s.Text()
		srv.MatchObservables(func(o common.Observable) bool {
			for field, regexp := range o.Regexps {
				switch field {
				case "message":
					if !regexp.MatchString(line) {
						return false
					}
				default:
					return false
				}
			}
			return true
		})
	}
	if err := s.Err(); err != nil {
		log.Fatalf("Failed to read from stdin: %v", err)
	}
}

func jsonConsumer() {
	d := json.NewDecoder(os.Stdin)
	srv.WaitForInitialConnection()
	for {
		var m map[string]any
		if err := d.Decode(&m); err != nil {
			log.Fatalf("Failed to read JSON from stdin: %v", err)
		}
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
