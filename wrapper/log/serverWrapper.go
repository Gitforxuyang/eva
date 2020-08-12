package log

import (
	"context"
	"fmt"
	config2 "github.com/Gitforxuyang/eva/config"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/utils"
	"google.golang.org/grpc"
	"time"
)

func NewServerWrapper() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := logger.GetLogger()
	config := config2.GetConfig()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			emap := map[string]interface{}{}
			if err != nil {
				emap = utils.StructToMap(error2.DecodeStatus(err))
			}
			if config.GetLogConfig().Server {
				log.Info(ctx, "收到的请求", logger.Fields{
					"req":     utils.StructToMap(req),
					"resp":    utils.StructToMap(resp),
					"method":  info.FullMethod,
					"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
					"err":     emap,
				})
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}
}
