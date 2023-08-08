logwait_common.pb.go: logwait_common.proto observer.proto scanner.proto Makefile
	protoc \
		--go_opt=paths=source_relative                                  \
		--go_out=.                                                      \
		--go-grpc_out=.                                                 \
		--go-grpc_opt=paths=source_relative                             \
		--go-grpc_opt=require_unimplemented_servers=false               \
		logwait_common.proto observer.proto scanner.proto
