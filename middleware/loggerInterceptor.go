package middleware

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func Logger() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return logger
}

func logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	defer func() {
		logrus.WithFields(logrus.Fields{
			"req":     req,
			"resp":    resp,
			"method":  info.FullMethod,
			"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
		}).Info("收到的请求")
	}()
	return handler(ctx, req)
}
