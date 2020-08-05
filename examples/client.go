package main

import (
	"context"
	client2 "github.com/Gitforxuyang/eva/client"
	"github.com/Gitforxuyang/eva/proto"
	"github.com/Gitforxuyang/eva/util/logger"
)

func main() {
	logger.Init("demo client")
	client := client2.GetGRpcSayHelloServiceClient()
	client.Hello(context.TODO(), &hello.String{})

}
