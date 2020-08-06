package utils

import "context"

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
