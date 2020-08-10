package server

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	trace2 "github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/Gitforxuyang/eva/wrapper/catch"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/Gitforxuyang/eva/wrapper/trace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type RegisterService func(server *grpc.Server)

//注册关闭服务时的回调
type RegisterShutdown func()

var (
	grpcServer   *grpc.Server
	listen       net.Listener
	shutdownFunc []RegisterShutdown = make([]RegisterShutdown, 0)
)

func Init(registerService RegisterService) {
	config.Init()
	conf := config.GetConfig()
	var err error
	listen, err = net.Listen("tcp", fmt.Sprintf(":%d", conf.GetPort()))
	utils.Must(err)
	logger.Init(conf.GetName())
	trace2.Init(fmt.Sprintf("%s_%s", conf.GetName(), conf.GetENV()),
		conf.GetTraceConfig().Endpoint, conf.GetTraceConfig().Ratio)
	grpcServer = grpc.NewServer(
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

}

func Run() {
	go func() {
		err := grpcServer.Serve(listen)
		utils.Must(err)
	}()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL)
	s := <-c
	logger.GetLogger().Info(context.TODO(), "signal", logger.Fields{
		"signal": s.String(),
	})
	grpcServer.GracefulStop()
	//做一些资源关闭动作
	for _, v := range shutdownFunc {
		v()
	}
}

func RegisterShutdownFunc(shutdown RegisterShutdown) {
	shutdownFunc = append(shutdownFunc, shutdown)
}
