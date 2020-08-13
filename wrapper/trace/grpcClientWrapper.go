package trace

import (
	"context"
	"github.com/Gitforxuyang/eva/config"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
)

func NewClientWrapper(tracer *trace.Tracer) func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	conf := config.GetConfig().GetTraceConfig()
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if conf.GRpcClient {
			ctx, span, err := tracer.StartGRpcClientSpanFromContext(ctx, method)
			if err != nil {
				logger.GetLogger().Error(ctx, "链路错误", logger.Fields{"err": utils.StructToMap(err)})
			}
			defer span.Finish()
			err = invoker(ctx, method, req, reply, cc, opts...)
			if conf.Log {
				span.LogFields(
					log.Object("req", utils.StructToJson(req)),
					log.Object("resp", utils.StructToJson(reply)),
				)
			}
			if err != nil {
				ext.Error.Set(span, true)
				span.LogFields(log.String("event", "error"))
				span.LogFields(
					log.Object("evaError", utils.StructToJson(err)),
				)
			}
			return err
		} else {
			err := invoker(ctx, method, req, reply, cc, opts...)
			return err
		}
	}
}
