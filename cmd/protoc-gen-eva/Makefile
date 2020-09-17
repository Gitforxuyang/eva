

proto:
	./examples/*.eva.go || true
	protoc -I=./examples/googleapis   -I=./examples --eva_out=plugins=all:./examples hello.proto


grpc:
	protoc  -I=./examples/googleapis -I=./examples --go_out=plugins=grpc:./examples hello.proto

run:
	go install github.com/Gitforxuyang/eva/cmd/protoc-gen-eva
	make proto