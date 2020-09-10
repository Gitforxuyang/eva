package main

import (
	"fmt"
	"github.com/Gitforxuyang/eva/examples/conf"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/examples/service"
	mongo2 "github.com/Gitforxuyang/eva/plugin/mongo"
	"github.com/Gitforxuyang/eva/plugin/redis"
	"github.com/Gitforxuyang/eva/server"
	"github.com/Gitforxuyang/eva/util/utils"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println(utils.StructToJson(hello.GetServerDesc()))
	server.Init()
	conf.Registry()
	rdb := redis.GetRedisClient("node")
	mongo := mongo2.GetMongoClient("node")
	server.RegisterGRpcService(func(server *grpc.Server) {
		//reflection.Register(server)
		hello.RegisterSayHelloServiceServer(server, service.NewHelloServiceServer(rdb, mongo))
	},hello.GetServerDesc())
	//client:=client2.GetGRpcSayHelloServiceClient()
	server.Run()

	//listen, err := net.Listen("tcp", ":50001")
	//if err != nil {
	//	panic(err)
	//}
	//logger.Init("demo")
	//tracer, err := trace2.NewTracer("eva_local", "http://192.168.3.23:14268/api/traces")
	//grpcServer := grpc.NewServer(
	//	grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
	//		trace.NewGRpcServerWrapper(tracer),
	//		log.NewServerWrapper(),
	//		catch.NewServerWrapper(),
	//	)),
	//	grpc.KeepaliveParams(keepalive.ServerParameters{
	//		MaxConnectionIdle: time.Second * 50,
	//		//MaxConnectionAge:time.Second*20,
	//	}),
	//)
	//hello.RegisterSayHelloServiceServer(grpcServer, &service.HelloServiceServer{})
	//err = grpcServer.Serve(listen)
	//if err != nil {
	//	panic(err)
	//}
}
