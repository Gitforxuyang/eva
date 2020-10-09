
proto:
	mkdir ./proto/example || true
	protoc --eva_out=plugins=all:./proto/example -I=./proto example.proto
	protoc --go_out=plugins=grpc:./proto/example -I=./proto example.proto

.PHONY: proto

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o example main.go

docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o example main.go
	docker build -t eva-demo:v1 .
