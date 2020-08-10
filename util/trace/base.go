package trace

import (
	"github.com/Gitforxuyang/eva/util/utils"
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

var (
	tracer *Tracer
)

func GetTracer() *Tracer {
	utils.NotNil(tracer, "tracer")
	return tracer

}

func Init(serviceName string, addr string, ratio float64) {
	t, err := newTracer(serviceName, addr, ratio)
	utils.Must(err)
	tracer = t
}

func newTracer(serviceName string, addr string, ratio float64) (*Tracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: ratio,
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
