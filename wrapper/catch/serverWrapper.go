package catch

import (
	"context"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/getsentry/sentry-go"
	"github.com/go-errors/errors"
	"github.com/opentracing/opentracing-go"
	log2 "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func NewServerWrapper() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := logger.GetLogger()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = error2.PanicError
				log.Error(ctx, "发生panic", logger.Fields{"e": e})
				span, ok := ctx.Value("_span").(opentracing.Span)
				if ok {
					span.LogFields(log2.Object("stack", zap.Stack("stack")))
				}
				sentry.CaptureException(errors.New(e))
			}
		}()
		deadline, _ := ctx.Deadline()
		//如果超时5s在deadline之后，则重置deadline为5s后
		if time.Now().Add(time.Second * 5).After(deadline) {
			ctx, _ = context.WithTimeout(ctx, time.Second*5)
		}
		resp, err = handler(ctx, req)
		if err != nil {
			evaError := error2.FromError(err)
			return resp, error2.EncodeStatus(evaError).Err()
		}
		return resp, err
	}
}
