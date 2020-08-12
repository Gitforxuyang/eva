package service

import (
	"context"
	"github.com/Gitforxuyang/eva/examples/proto/hello"
	"github.com/Gitforxuyang/eva/plugin/mongo"
	"github.com/Gitforxuyang/eva/plugin/redis"
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
	a := &Animal{String: "str"}
	_, err := m.mongo.Database("demo").Collection("demo").InsertOne(ctx, a)
	if err != nil {
		return nil, err
	}
	return &hello.String{}, nil
}
