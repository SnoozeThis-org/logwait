snoozethis.pb.go: snoozethis.proto Makefile
	protoc \
		--go_opt=paths=source_relative                                  \
		--go_out=.                                                      \
		--go-grpc_out=.                                                 \
		--go-grpc_opt=paths=source_relative                             \
		--go-grpc_opt=require_unimplemented_servers=false               \
		snoozethis.proto
