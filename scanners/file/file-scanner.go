package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/SnoozeThis-org/logwait/scanners/common"
	"google.golang.org/grpc"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}
	defer c.Close()

	srv := common.NewService(c)

	go srv.ConnectLoop()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatalf("Cannot open file %q: %v", flag.Arg(0), err)
	}

	defer file.Close()

	// TODO: Add support for truncated files
	// TODO: Add support for tailing multiple files
	// TODO: Stop tailing when no observables?

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				time.Sleep(20 * time.Millisecond)
				continue
			}
			break
		}

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
}
