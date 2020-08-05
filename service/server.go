package service

import (
	"context"
	hello "github.com/Gitforxuyang/eva/proto"
)

type HelloServiceServer struct {
}

func (HelloServiceServer) Hello(context.Context, *hello.String) (*hello.String, error) {
	return &hello.String{}, nil
}
