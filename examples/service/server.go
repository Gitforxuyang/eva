package service

import (
	"context"
	"errors"
	hello "github.com/Gitforxuyang/eva/examples/proto"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"google.golang.org/grpc/codes"
	"time"
)

type HelloServiceServer struct {
}

func (HelloServiceServer) Hello(ctx context.Context, req *hello.String) (*hello.String, error) {
	time.Sleep(time.Second * 4)
	select {
	case <-ctx.Done():
		return nil, errors.New("context done")
	default:

	}
	return &hello.String{}, error2.New("demo", "自定义错误", "测试用的", 10001, codes.Internal)
}
