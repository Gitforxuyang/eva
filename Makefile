

eva:
	protoc -I=./examples/proto --go_out=plugins=grpc:./examples/proto/hello hello.proto
	protoc -I=./examples/proto --eva_out=plugins=all:./examples/proto/hello hello.proto