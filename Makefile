

eva:
	protoc -I=./proto --go_out=plugins=grpc:./proto hello.proto
	protoc -I=./proto --eva_out=plugins=all:./proto hello.proto