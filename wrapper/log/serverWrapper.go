package log

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/utils"
	"google.golang.org/grpc"
	"time"
)

func NewServerWrapper() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := logger.GetLogger()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			log.Info(ctx, "收到的请求", logger.Fields{
				"req":     utils.StructToMap(req),
				"resp":    utils.StructToMap(resp),
				"method":  info.FullMethod,
				"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
				"err":     utils.StructToMap(err),
			})
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
