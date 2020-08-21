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
	tracer  opentracing.Tracer
	closer  io.Closer
	sampler *jaeger.ProbabilisticSampler
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

func SetRatio(ratio float64) {
	tracer.sampler.Update(ratio)
}
func newTracer(serviceName string, addr string, ratio float64) (*Tracer, error) {
	cfg := config.Configuration{
		ServiceName: serviceName,
		//Sampler: &config.SamplerConfig{
		//	Type:  jaeger.SamplerTypeProbabilistic,
		//	Param: ratio,
		//},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	sampler, _ := jaeger.NewProbabilisticSampler(1)
	sender := transport.NewHTTPTransport(addr)
	reporter := jaeger.NewRemoteReporter(sender)

	tracer, closer, err := cfg.NewTracer(
		config.Reporter(reporter),
		config.Sampler(sampler),
	)
	if err != nil {
		return nil, err
	}
	t := &Tracer{
		tracer:  tracer,
		closer:  closer,
		sampler: sampler,
	}
	return t, nil
}
