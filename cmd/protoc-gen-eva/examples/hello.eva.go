package hello

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/client/selector"
	"github.com/Gitforxuyang/eva/config"
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
		grpc.WithBlock(),
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
