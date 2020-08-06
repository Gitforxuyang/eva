package catch

import (
	"context"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"google.golang.org/grpc"
	"time"
)

func NewServerWrapper() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log := logger.GetLogger()
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := recover(); e != nil {
				err = error2.Parse("发生异常")
				log.Panic(ctx, "发生panic", logger.Fields{"e": e})
				//TODO:sentry捕获
			}
		}()
		deadline, _ := ctx.Deadline()
		//如果超时5s在deadline之后，则重置deadline为5s后
		if time.Now().Add(time.Second * 3).After(deadline) {
			ctx, _ = context.WithTimeout(ctx, time.Second*3)
		}
		resp, err = handler(ctx, req)
		if err != nil {
			evaError := error2.FromError(err)
			return resp, error2.EncodeStatus(evaError).Err()
		}
		return resp, err
	}
}
