module github.com/SnoozeThis-org/logwait

go 1.18

require (
	github.com/Jille/convreq v1.6.1-0.20230808184511-45e4d7abdccd
	github.com/Jille/easymutex v1.1.0
	github.com/Jille/genericz v0.4.0
	github.com/fsnotify/fsnotify v1.6.0
	github.com/grafana/loki/pkg/push v0.0.0-20230919104151-a8d5815510bd
	github.com/influxdata/go-syslog/v3 v3.0.0
	github.com/klauspost/compress v1.17.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
)

replace github.com/grafana/loki/pkg/push v0.0.0-20230919104151-a8d5815510bd => github.com/Jille/loki/pkg/push v0.0.0-20230919191926-b7b103c7c7f3
