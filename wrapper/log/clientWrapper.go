package log

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/util/logger"
	"google.golang.org/grpc"
	"time"
)

func NewClientWrapper() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log := logger.GetLogger()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		defer func() {
			log.Info(ctx, "发起的请求", logger.Fields{
				"req":     req,
				"reply":   reply,
				"method":  method,
				"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
			})
		}()
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}
