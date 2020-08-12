package service

import (
	"context"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/plugin/redis"
	"time"
)

type HelloServiceServer struct {
	redis redis.EvaRedis
}

func NewHelloServiceServer(rdb redis.EvaRedis) *HelloServiceServer {
	return &HelloServiceServer{redis: rdb}
}

func (m *HelloServiceServer) Hello(ctx context.Context, req *hello.String) (*hello.String, error) {
	ctx, _ = context.WithTimeout(ctx, time.Second*1000)
	_, err := m.redis.Set(ctx, "rdb:demo", "1", 0)
	if err != nil {
		return nil, err
	}
	return &hello.String{}, nil
}
