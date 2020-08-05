package main

import (
	"context"
	client2 "github.com/Gitforxuyang/eva/client"
	"github.com/Gitforxuyang/eva/proto"
)

func main() {
	client := client2.GetGRpcSayHelloServiceClient()
	client.Hello(context.TODO(), &hello.String{})

}
