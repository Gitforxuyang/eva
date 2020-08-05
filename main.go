package main

import (
	"github.com/Gitforxuyang/eva/proto"
	"github.com/Gitforxuyang/eva/service"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/wrapper/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
)

func main() {
	listen, err := net.Listen("tcp", ":50001")
	if err != nil {
		panic(err)
	}
	//logrus.SetFormatter(&logrus.JSONFormatter{})
	logger.Init("demo")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			//middleware.Logger(),
			log.NewServerWrapper(),
		)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Second * 50,
			//MaxConnectionAge:time.Second*20,
		}),
	)
	hello.RegisterSayHelloServiceServer(grpcServer, &service.HelloServiceServer{})
	err = grpcServer.Serve(listen)
	if err != nil {
		panic(err)
	}
}
