package main

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/SnoozeThis-org/logwait/config"
	"github.com/SnoozeThis-org/logwait/scanners/common"
	syslog "github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/octetcounting"
	"github.com/influxdata/go-syslog/v3/rfc3164"
	"github.com/influxdata/go-syslog/v3/rfc5424"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var (
	listenTCP = pflag.String("tcp", "", "Address to listen for syslog messages (TCP), for example :514")
	listenUDP = pflag.String("udp", "", "Address to listen for syslog messages (UDP)")
	fmt3164   = pflag.Bool("rfc3164", false, "Syslog messages confirm to RFC3164")
	fmt5424   = pflag.Bool("rfc5424", false, "Syslog messages confirm to RFC5424")

	srv     *common.Service
	learner *common.FieldLearner
)

func main() {
	config.Parse()

	if *fmt3164 && *fmt5424 {
		log.Fatalf("Cannot use both RFC3164 and RFC424")
	}
	if !*fmt3164 && !*fmt5424 {
		log.Fatalf("Please select RFC3164 or RFC424")
	}
	if *fmt3164 && *listenTCP != "" {
		log.Fatalf("Cannot use RFC3164 with TCP")
	}

	c, err := grpc.Dial(*common.ObserverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Cannot connect to observer at %q: %v", *common.ObserverAddress, err)
	}

	srv = common.NewService(c, "syslog-scanner")
	learner = srv.CreateFieldLearner('.')
	learner.SetStaticFields([]string{"message", "appname", "hostname", "procid", "msgid", "facility", "severity"})

	if *listenTCP != "" {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *listenTCP)
		if err != nil {
			log.Fatalf("Invalid TCP port %v: %v", *listenTCP, err)
		}
		l, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			log.Fatalf("Cannot open TCP port %v: %v", *listenTCP, err)
		}
		go handleTCP(l)
	}

	if *listenUDP != "" {
		udpAddr, err := net.ResolveUDPAddr("udp", *listenUDP)
		if err != nil {
			log.Fatalf("Invalid UDP port %v: %v", *listenUDP, err)
		}
		l, err := net.ListenUDP("udp", udpAddr)
		if err != nil {
			log.Fatalf("Cannot open UDP port %v: %v", *listenUDP, err)
		}
		go handleUDP(l)
	}

	srv.ConnectLoop()
}

func handleTCP(l *net.TCPListener) {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		go func() {
			defer conn.Close()

			acc := func(res *syslog.Result) {
				if res.Error != nil {
					log.Printf("Failed to parse log message: %v", err)
					return
				}
				newMessage(res.Message)
			}
			octetcounting.NewParser(syslog.WithBestEffort(), syslog.WithListener(acc)).Parse(conn)
		}()
	}
}

func handleUDP(l *net.UDPConn) {
	var p syslog.Machine
	if *fmt3164 {
		p = rfc3164.NewParser(rfc3164.WithBestEffort())
	} else {
		p = rfc5424.NewParser(rfc5424.WithBestEffort())
	}

	resp := make([]byte, 9000)
	for {
		n, err := l.Read(resp)
		if err == nil {
			msg, err := p.Parse(resp[:n])
			if err != nil {
				log.Printf("Failed to parse log message: %v", err)
				continue
			}
			newMessage(msg)
		}
	}
}

func newMessage(msg syslog.Message) {
	var base syslog.Base
	var structuredData map[string]map[string]string
	switch v := msg.(type) {
	case *rfc5424.SyslogMessage:
		base = v.Base
		if v.StructuredData != nil {
			structuredData = *v.StructuredData
			for k1, m := range structuredData {
				for k2 := range m {
					learner.Seen(k1 + "." + k2)
				}
			}
		}
	case *rfc3164.SyslogMessage:
		base = v.Base
	default:
		log.Printf("Unspported log type: %T", msg)
		return
	}

	srv.MatchObservables(func(o common.Observable) bool {
		for field, regexp := range o.Regexps {
			var v *string
			switch field {
			case "message":
				v = base.Message
			case "appname":
				v = base.Appname
			case "hostname":
				v = base.Hostname
			case "procid":
				v = base.ProcID
			case "msgid":
				v = base.MsgID
			case "facility":
				v = base.FacilityLevel()
			case "severity":
				v = base.SeverityLevel()
			default:
				sp := strings.SplitN(field, ".", 2)
				if len(sp) != 2 {
					return false
				}
				s, ok := structuredData[sp[0]][sp[1]]
				if !ok {
					return false
				}
				v = &s
			}
			if v == nil {
				return false
			}
			if !regexp.MatchString(*v) {
				return false
			}
		}
		return true
	})
}
