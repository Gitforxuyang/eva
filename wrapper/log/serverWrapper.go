package log

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/util/logger"
	"google.golang.org/grpc"
	"time"
)

func NewServerWrapper() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := logger.GetLogger()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			log.Info(ctx, "收到的请求", logger.Fields{
				"req":     req,
				"resp":    resp,
				"method":  info.FullMethod,
				"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
			})
		}()
		return handler(ctx, req)
	}
}
