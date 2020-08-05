package utils

import "context"

func GetTraceId(ctx context.Context) string {
	traceId := ctx.Value("traceId")
	if traceId == nil {
		return "nil"
	} else {
		return traceId.(string)
	}
}
