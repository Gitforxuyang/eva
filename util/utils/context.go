package utils

import (
	"context"
	"fmt"
	error2 "github.com/Gitforxuyang/eva/util/error"
)

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ctx.Value("traceId")
	if traceId == nil {
		return ""
	} else {
		return traceId.(string)
	}
}

func ContextDie(ctx context.Context) error {
	select {
	case <-ctx.Done():
		fmt.Print("ContextDieError")
		return error2.ContextDieError
	default:
		return nil
	}
}
