package middleware

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

func ClientLogger(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	defer func() {
		logrus.WithFields(logrus.Fields{
			"req":     req,
			"reply":   reply,
			"method":  method,
			"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
		}).Info("发起的请求")
	}()
	err := invoker(ctx, method, req, reply, cc, opts...)
	return err
}
