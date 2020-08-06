package service

import (
	"context"
	hello "github.com/Gitforxuyang/eva/proto"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"google.golang.org/grpc/codes"
)

type HelloServiceServer struct {
}

func (HelloServiceServer) Hello(context.Context, *hello.String) (*hello.String, error) {
	return &hello.String{}, error2.New("demo", "自定义错误", "测试用的", 10001, codes.Internal)
}