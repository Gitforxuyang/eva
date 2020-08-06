package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/transport"
	"io"
	"time"
)

type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
}

func NewTracer(serviceName string, addr string) (*Tracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	sender := transport.NewHTTPTransport(addr)
	reporter := jaeger.NewRemoteReporter(sender)
	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
	)
	if err != nil {
		return nil, err
	}
	t := &Tracer{
		tracer: tracer,
		closer: closer,
	}
	return t, nil
}
