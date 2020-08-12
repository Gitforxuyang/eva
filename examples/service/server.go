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
	lock, err := m.redis.Lock(ctx, "demo:lock:key", time.Second*20, redis.LockOptions(time.Millisecond*0, 0))
	if err != nil {
		return nil, err
	}
	err = lock.UnLock(ctx)
	if err != nil {
		return nil, err
	}
	return &hello.String{}, nil
}
