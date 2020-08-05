package client

import (
	"context"
	"github.com/Gitforxuyang/eva/client/selector"
	"github.com/Gitforxuyang/eva/proto"
	"github.com/Gitforxuyang/eva/wrapper/log"
	"github.com/sirupsen/logrus"
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
	conn, err := grpc.Dial(":50001",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithBalancerName(roundrobin.Name),
		grpc.WithResolvers(selector.NewCustomResolverBuilder("dns")),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second * 10,
				Timeout:             time.Second * 1,
				PermitWithoutStream: true,
			}),
		grpc.WithChainUnaryInterceptor(log.NewClientWrapper()),
	)
	c := &grpcSayHelloServiceClient{}
	c.client = hello.NewSayHelloServiceClient(conn)
	if err != nil {
		panic(err)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return c
}
