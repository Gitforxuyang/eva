package trace

import (
	"context"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
)

func NewGrpcServerWrapper(tracer *trace.Tracer) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, span, err := tracer.StartSpanFromContext(ctx, info.FullMethod)
		span.SetTag("span.kind", "server")
		if err != nil {
			logger.GetLogger().Error(ctx, "链路错误", logger.Fields{"err": err})
		}
		defer span.Finish()
		resp, err = handler(ctx, req)
		span.LogFields(
			log.Object("req", utils.StructToJson(req)),
			log.Object("resp", utils.StructToJson(resp)),
		)
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"))
			span.LogFields(
				log.Object("evaError", utils.StructToJson(err)),
			)
		}
		return resp, err
	}
}
