package server

import (
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	trace2 "github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/wrapper/catch"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/Gitforxuyang/eva/wrapper/trace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"net"
)

type RegisterService func(server *grpc.Server)

func Run(registerService RegisterService) {
	config.Init()
	conf := config.GetConfig()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GetPort()))
	if err != nil {
		panic(err)
	}
	logger.Init(conf.GetName())
	trace2.Init(fmt.Sprintf("%s_%s", conf.GetName(), conf.GetENV()),
		conf.GetTraceConfig().Endpoint, conf.GetTraceConfig().Ratio)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			trace.NewGRpcServerWrapper(trace2.GetTracer()),
			log.NewServerWrapper(),
			catch.NewServerWrapper(),
		)),
		//grpc.KeepaliveParams(keepalive.ServerParameters{
		//MaxConnectionIdle: time.Second * 50,
		//MaxConnectionAge:time.Second*20,
		//}),
	)
	registerService(grpcServer)
	//hello.RegisterSayHelloServiceServer(grpcServer, &service.HelloServiceServer{})
	err = grpcServer.Serve(listen)
	if err != nil {
		panic(err)
	}
}
