package main

import (
	"context"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"strconv"
	"time"
)

func main() {
	logger.Init("demo client")
	trace.Init("eva_local", "http://192.168.3.23:14268/api/trace", 1)
	config.Init()
	client := hello.GetGRpcSayHelloServiceClient()
	client.Hello(context.TODO(), &hello.String{Name: strconv.Itoa(int(time.Now().Unix()))})
}
