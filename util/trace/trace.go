package trace

import (
	"context"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/grpc/metadata"
	"strings"

	"github.com/opentracing/opentracing-go"
)

// metadataReaderWriter satisfies both the opentracing.TextMapReader and
// opentracing.TextMapWriter interfaces.
type metadataReaderWriter struct {
	metadata.MD
}

func (w metadataReaderWriter) Set(key, val string) {
	// The GRPC HPACK implementation rejects any uppercase keys here.
	//
	// As such, since the HTTP_HEADERS format is case-insensitive anyway, we
	// blindly lowercase the key (which is guaranteed to work in the
	// Inject/Extract sense per the OpenTracing spec).
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

func (w metadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Tracer) StartGRpcServerSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = make(map[string][]string)
	}
	md = md.Copy()
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := m.tracer.Extract(opentracing.TextMap, metadataReaderWriter{md}); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	sp := m.tracer.StartSpan(name, opts...)
	ctx = withTraceId(ctx, sp)
	ext.SpanKindRPCServer.Set(sp)
	ext.Component.Set(sp, "grpc")
	ctx = opentracing.ContextWithSpan(ctx, sp)
	return ctx, sp, nil
}

func (m *Tracer) StartGRpcClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	//ctx = withTraceId(ctx, sp)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "grpc")
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := metadataReaderWriter{md}
	err := m.tracer.Inject(sp.Context(), opentracing.TextMap, mdWriter)
	if err != nil {
		return nil, nil, err
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, sp, nil
}

func (m *Tracer) StartHttpClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	//ctx = withTraceId(ctx, sp)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "http")
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := metadataReaderWriter{md}
	err := m.tracer.Inject(sp.Context(), opentracing.HTTPHeaders, mdWriter)
	if err != nil {
		return nil, nil, err
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, sp, nil
}

func (m *Tracer) StartRedisClientSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	}
	sp := m.tracer.StartSpan(name, opts...)
	ext.SpanKindRPCClient.Set(sp)
	ext.Component.Set(sp, "redis")
	ext.DBType.Set(sp, "redis")
	return ctx, sp, nil
}

func withTraceId(ctx context.Context, span opentracing.Span) context.Context {
	s, ok := span.Context().(jaeger.SpanContext)
	if ok {
		ctx = context.WithValue(ctx, "traceId", s.TraceID().String())
	}
	return ctx
}
