package catch

import (
	"context"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"google.golang.org/grpc"
)

//用来将其它服务的返回错误转换为eva定义的错规范
func NewClientWrapper() func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			e := error2.DecodeStatus(err)
			err = e
		}
		return err
	}
}
