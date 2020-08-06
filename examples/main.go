package main

import (
	hello "github.com/Gitforxuyang/eva/examples/proto"
	"github.com/Gitforxuyang/eva/examples/service"
	"github.com/Gitforxuyang/eva/util/logger"
	trace2 "github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/wrapper/catch"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/Gitforxuyang/eva/wrapper/trace"
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
	logger.Init("demo")
	tracer, err := trace2.NewTracer("eva_local", "http://192.168.3.23:14268/api/traces")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			trace.NewGRpcServerWrapper(tracer),
			log.NewServerWrapper(),
			catch.NewServerWrapper(),
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
