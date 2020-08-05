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
	start := time.Now().Unix()
	defer func() {
		logrus.WithFields(logrus.Fields{
			"req":     req,
			"resp":    resp,
			"method":  info.FullMethod,
			"useTime": fmt.Sprintf("%d s", time.Now().Unix()-start),
		}).Info("msg")
	}()
	return handler(ctx, req)
}
