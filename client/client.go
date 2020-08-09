package client

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/client/selector"
	"github.com/Gitforxuyang/eva/config"
	hello "github.com/Gitforxuyang/eva/examples/proto"
	trace2 "github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/wrapper/catch"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/Gitforxuyang/eva/wrapper/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	"time"
)

type GRpcSayHelloServiceClient interface {
	Hello(ctx context.Context, req *hello.String) (resp *hello.String, err error)
}

type grpcSayHelloServiceClient struct {
	client hello.SayHelloServiceClient
}

func (m *grpcSayHelloServiceClient) Hello(ctx context.Context, req *hello.String) (resp *hello.String, err error) {
	resp, err = m.client.Hello(ctx, req)
	return resp, err
}

func GetGRpcSayHelloServiceClient() GRpcSayHelloServiceClient {
	tracer := trace2.GetTracer()
	grpcClientConfig := config.GetConfig().GetGRpc("client")
	conn, err := grpc.Dial(fmt.Sprintf("%s", grpcClientConfig.Endpoint),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithResolvers(selector.NewCustomResolverBuilder(grpcClientConfig.Mode)),
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
	c.client = hello.NewSayHelloServiceClient(conn)
	if err != nil {
		panic(err)
	}
	return c
}
