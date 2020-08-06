package main

import (
	"context"
	client2 "github.com/Gitforxuyang/eva/client"
	"github.com/Gitforxuyang/eva/proto"
	"github.com/Gitforxuyang/eva/util/logger"
	"strconv"
	"time"
)

func main() {
	logger.Init("demo client")
	client := client2.GetGRpcSayHelloServiceClient()
	client.Hello(context.TODO(), &hello.String{Name: strconv.Itoa(int(time.Now().Unix()))})

}
