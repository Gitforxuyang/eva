// Code generated by protoc-gen-eva. DO NOT EDIT.
// source: hello.proto

package hello

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/client/selector"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/registory/etcd"
	"github.com/Gitforxuyang/eva/server"
	trace2 "github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/Gitforxuyang/eva/wrapper/catch"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/Gitforxuyang/eva/wrapper/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

//获取服务的描述信息
func GetServerDesc() *etcd.Service {
	messageMap := make(map[string]map[string]string, 1)
	messageString := make(map[string]string)
	messageString["name"] = "TYPE_STRING"
	messageString["a"] = ".hello.String.AEntry"
	messageString["i32"] = "TYPE_INT32"
	messageString["i64"] = "TYPE_INT64"
	messageString["f32"] = "TYPE_FLOAT"
	messageString["f64"] = "TYPE_DOUBLE"
	messageMap["String"] = messageString
	service := new(etcd.Service)
	service.Name = "SayHelloService"
	service.Package = "sayHelloService"
	service.AppId = "sayHelloService"
	service.Methods = make(map[string]etcd.Method, 1)
	methodHello := etcd.Method{}
	methodHello.Req = messageMap["String"]
	methodHello.Resp = messageMap["String"]
	service.Methods["Hello"] = methodHello
	return service
}

type GRpcSayHelloServiceClient interface {
	Hello(ctx context.Context, req *String) (resp *String, err error)
}
type grpcSayHelloServiceClient struct {
	client SayHelloServiceClient
}

func (m *grpcSayHelloServiceClient) Hello(ctx context.Context, req *String) (resp *String, err error) {
	resp, err = m.client.Hello(ctx, req)
	return resp, err
}
func GetGRpcSayHelloServiceClient() GRpcSayHelloServiceClient {
	tracer := trace2.GetTracer()
	grpcClientConfig := config.GetConfig().GetGRpc("SayHelloService")
	conn, err := grpc.Dial(fmt.Sprintf("%s", grpcClientConfig.Endpoint),
		grpc.WithInsecure(),
		grpc.WithBalancer(grpc.RoundRobin(selector.NewCustomResolverBuilder(grpcClientConfig.Mode).GetResolver(grpcClientConfig.Endpoint))),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second * 10,
				Timeout:             time.Second * 1,
				PermitWithoutStream: true,
			}),
		grpc.WithChainUnaryInterceptor(
			trace.NewClientWrapper(tracer),
			log.NewClientWrapper(),
			catch.NewClientWrapper(grpcClientConfig.Timeout),
		),
	)
	c := &grpcSayHelloServiceClient{}
	c.client = NewSayHelloServiceClient(conn)
	utils.Must(err)
	server.RegisterShutdownFunc(func() {
		conn.Close()
	})
	return c
}