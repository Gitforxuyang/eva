package log

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/utils"
	"google.golang.org/grpc"
	"time"
)

func NewClientWrapper() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log := logger.GetLogger()
	conf := config.GetConfig()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		var err error
		defer func() {
			if conf.GetLogConfig().GRpcClient {
				log.Info(ctx, "发起的请求", logger.Fields{
					"req":     utils.StructToMap(req),
					"resp":    utils.StructToMap(reply),
					"method":  method,
					"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
					"err":     utils.StructToMap(err),
				})
			}
		}()
		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
