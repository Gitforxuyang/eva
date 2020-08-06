package trace

import (
	"context"
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

func (m *Tracer) StartSpanFromContext(ctx context.Context, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
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

	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, metadataReaderWriter{md}); err != nil {
		return nil, nil, err
	}
	ctx = opentracing.ContextWithSpan(ctx, sp)
	//ctx = metadata.NewContext(ctx, md)
	return ctx, sp, nil
}
