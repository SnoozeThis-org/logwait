# Monitoring log files with SnoozeThis
logwait integrates your logging system with [SnoozeThis](https://www.snoozethis.com/), allowing you to snooze issues until a specific log message occurs.

A common scenario would be: you get assigned an bug report for one of your applications but find it difficult to reproduce. You add additional logging that would help you if the issue reoccurs. You then create a filter looking for a specific log message and hand over the issue to SnoozeThis, who will start monitoring your log files. As soon as a match has been found, SnoozeThis will assign the issue back to you, so you can start fixing the bug.

logwait consists of two parts: an observer and one or more scanners. Both run locally, so no log messages will ever leave your environment.

# Observer
The observer serves a simple web interface where you can create new observables using filters. The web interface will give you a URL which you can use to snooze an issue in your issue tracker by simply commenting <code>@SnoozeThis url</code>.
When a scanner (see below) detects a log message matching your filters, the observer will ping SnoozeThis, who in turn will reassing the issue back to you.

In order to use the observer you will need to create a token on our website. This token is linked to your SnoozeThis account (user or organization) and will make sure only you (or your coworkers) can setup snoozes for your logging application.

# Scanner(s)
A scanner is a simple program that digests log messages and scans these messages for the observables defined in the observer. We have scanners for syslog and tailing (text) log files, with more to come.

You can start as many scanners as you need and they can run on different hosts than the observer. The scanner connects to the observer using GRPC (by default on tcp port 1600).

# Creating your own scanner
If you have a need for a specific scanner for your logging application feel free to create an issue. You can also create your own scanner. Have a look at one of the existing scanners or start using this Go code:
```
package main

import (
	"flag"
	"log"

	"github.com/SnoozeThis-org/logwait/scanners/common"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

  // Connect to the observer
	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}
	defer c.Close()

	srv := common.NewService(c)
	go srv.ConnectLoop()

	for {
		// Digest events from your logging application
		event := `this is a log message`

		// Look for a match with one or more observables
		srv.MatchObservables(func(o common.Observable) bool {
			for field, regexp := range o.Regexps {
				switch field {
				case "message":
					if !regexp.MatchString(event) {
						return false
					}
				default:
					return false
				}
			}
      // Tell the observer a match has been found
			return true
		})
	}
}
```

The scanner communicates with the observer using streaming GRPC. If you prefer a different programming language for your scanner, have a look at the [proto](https://github.com/SnoozeThis-org/logwait/blob/master/proto/scanner.proto) file.
