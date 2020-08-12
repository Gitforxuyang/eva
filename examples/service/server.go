package service

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/plugin/mongo"
	"github.com/Gitforxuyang/eva/plugin/redis"
	"time"
)

type HelloServiceServer struct {
	redis redis.EvaRedis
	mongo mongo.EvaMongo
}

func NewHelloServiceServer(rdb redis.EvaRedis, mongo mongo.EvaMongo) *HelloServiceServer {
	return &HelloServiceServer{redis: rdb, mongo: mongo}
}

type Animal struct {
	String string `bson:"str"`
}

func (m *HelloServiceServer) Hello(ctx context.Context, req *hello.String) (*hello.String, error) {
	_, err := m.mongo.Database("demo").Collection("demo").InsertOne(ctx, a)
	if err != nil {
		fmt.Println(err)
	}
	return &hello.String{}, nil
}
