all: common.pb.go observer.pb.go scanner.pb.go

.PHONY: all

common.pb.go: common.proto Makefile
	protoc \
		--go_opt=paths=source_relative                                  \
		--go_out=.                                                      \
		--go-grpc_out=.                                                 \
		--go-grpc_opt=paths=source_relative                             \
		--go-grpc_opt=require_unimplemented_servers=false               \
		common.proto

observer.pb.go: observer.proto Makefile
	protoc \
		--go_opt=paths=source_relative                                  \
		--go_out=.                                                      \
		--go-grpc_out=.                                                 \
		--go-grpc_opt=paths=source_relative                             \
		--go-grpc_opt=require_unimplemented_servers=false               \
		observer.proto

scanner.pb.go: scanner.proto Makefile
	protoc \
		--go_opt=paths=source_relative                                  \
		--go_out=.                                                      \
		--go-grpc_out=.                                                 \
		--go-grpc_opt=paths=source_relative                             \
		--go-grpc_opt=require_unimplemented_servers=false               \
		scanner.proto
