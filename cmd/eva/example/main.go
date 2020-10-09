
package main

import (
	"example/conf"
	"example/handler"
	"example/proto/example"
	"github.com/Gitforxuyang/eva/server"
	"google.golang.org/grpc"
)

func main(){
	server.Init()
	conf.Registry()
	server.RegisterGRpcService(func(server *grpc.Server) {
		example.RegisterExampleServer(server,&handler.HandlerService{})
	},example.GetServerDesc())
	server.Run()
}
