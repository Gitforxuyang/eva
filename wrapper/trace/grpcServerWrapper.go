package trace

import (
	"context"
	"github.com/Gitforxuyang/eva/config"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
)

func NewGRpcServerWrapper(tracer *trace.Tracer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	conf := config.GetConfig().GetTraceConfig()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, span, err := tracer.StartGRpcServerSpanFromContext(ctx, info.FullMethod)
		if err != nil {
			logger.GetLogger().Error(ctx, "链路错误", logger.Fields{"err": utils.StructToMap(err)})
		}
		defer span.Finish()
		resp, err = handler(ctx, req)
		if conf.Log {
			span.LogFields(
				log.Object("req", utils.StructToJson(req)),
				log.Object("resp", utils.StructToJson(resp)),
			)
		}
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"))
			span.LogFields(
				log.Object("evaError", utils.StructToJson(error2.DecodeStatus(err))),
			)
		}
		return resp, err
	}
}
