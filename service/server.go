package service

import (
	"context"
	hello "github.com/Gitforxuyang/eva/proto"
	"time"
)

type HelloServiceServer struct {
}

func (HelloServiceServer) Hello(context.Context, *hello.String) (*hello.String, error) {
	time.Sleep(time.Second * 10)
	return &hello.String{}, nil
}
