# Monitoring log files with SnoozeThis
logwait integrates your logging system with [SnoozeThis](https://www.snoozethis.com/), allowing you to snooze issues until a specific log message occurs.

A common scenario would be: you get assigned an bug report for one of your applications but find it difficult to reproduce. You add additional logging that would help you if the issue reoccurs. You then create a filter looking for a specific log message and hand over the issue to SnoozeThis, who will start monitoring your log files. As soon as a match has been found, SnoozeThis will assign the issue back to you, so you can start fixing the bug.

logwait consists of two parts: an observer and one or more scanners. Both run locally, so no log messages will ever leave your servers.

# Observer
The observer serves a simple web interface where you can create new observables using filters. The web interface will give you a URL which you can use to snooze an issue in your issue tracker by simply commenting <code>@SnoozeThis url</code>.
When a scanner (see below) detects a log message matching your filters, the observer will ping SnoozeThis, which in turn will reassign the issue back to you.

In order to use the observer you will need to create a token on our website. This token is linked to your SnoozeThis account (user or organization) and will make sure only you (or your coworkers) can setup snoozes for your logging application.

# Scanner(s)
A scanner is a simple program that digests log messages and scans these messages for the observables defined in the observer. We have scanners for syslog and tailing (text) log files, with more to come.

You can start as many scanners as you need and they can run on different hosts than the observer. The scanner connects to the observer using gRPC (by default on tcp port 1600).

# Running the observer
1. Download the latest version at https://github.com/SnoozeThis-org/logwait/releases/latest
2. Get a token at https://www.snoozethis.com/logs/
3. Create a config file (or use command line arguments or environment variables):
   ```
   {
     "http-port": 8080,
     "grpc-port": 1600,
     "token": "token from step 2",
     "signing-key": "secret"
   }
   ```
   You can choose your own ports and signing key. Make sure your coworkers can access the UI at the http port and the scanner(s) can connect to the observer using the grpc port.
4. Start the observer
   ```
   observer --config /path/to/config.json
   ```
5. Start one or more scanners (see below)
6. Point your browser to the http port and create a new observable

# Running the syslog scanner
The syslog scanner accepts log messages in RFC3164 or RFC5424 format via TCP or UDP
1. Make sure your observer (see above) is running
2. Download the latest version at https://github.com/SnoozeThis-org/logwait/releases/latest
3. Start the scanner
   ```
   # For RFC3164 messages via UDP
   syslog-scanner -observer-address observer-ip:1600 -udp :514 -rfc3164
   # For RFC5424 messages via TCP
   syslog-scanner -observer-address observer-ip:1600 -tcp :514 -rfc5424
   ```
4. Have your syslog forward messages to the scanner. This depends on your version and flavour of syslog, but this will probably work:
   ```
   *.* @scanner-ip
   ```

# Running the file scanner
The file scanner tails one or more files
1. Make sure your observer (see above) is running
2. Download the latest version at https://github.com/SnoozeThis-org/logwait/releases/latest
3. Start the scanner
   ```
   file-scanner -observer-address observer-ip:1600 <file>
   ```

# Creating your own scanner
If you have a need for a specific scanner for your logging application feel free to create an [issue](https://github.com/SnoozeThis-org/logwait/issues/new). You can also create your own scanner. Have a look at one of the existing scanners or start using this Go code:
```
package main

import (
	"log"

	"github.com/SnoozeThis-org/logwait/config"
	"github.com/SnoozeThis-org/logwait/scanners/common"
	"google.golang.org/grpc"
)

func main() {
	config.Parse()

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

The scanner communicates with the observer using streaming gRPC. If you prefer a different programming language for your scanner, have a look at the [proto](https://github.com/SnoozeThis-org/logwait/blob/master/proto/scanner.proto) file.

# Security

We understand you are careful about your log lines. Surely *your company* isn't logging plain-text passwords, but it's still better to keep your logs secure. We have the following security principles to ensure no data leaks can happen based on your usage of logwait (and SnoozeThis).

1. No log line ever leaves your servers. The Scanner just tells the Observer "there was a match", which the observer relays to SnoozeThis.
2. Filters can only be created by people able to access your Observer webinterface. SnoozeThis tells your Observer which filters are active, but when your Observer created the URL for them, it signed them with an HMAC. Because of this SnoozeThis employees can't prod around your log lines by trying different filters, because they wouldn't have been signed by your Observer.
3. The URL created by the Observer to snooze on *does* contain all the filters, because it's convenient for your colleagues to see what exactly you're snoozing on.
