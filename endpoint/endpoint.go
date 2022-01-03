package endpoint

import (
	"add-svc/svc"
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/log"
	stdopentracing "github.com/opentracing/opentracing-go"
	stdzipkin "github.com/openzipkin/zipkin-go"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	AddEndpoint endpoint.Endpoint
}

type AddRequest struct {
	A int32
	B int32
}

type AddResponse struct {
	V   int32
	Err error
}

func (res AddResponse) Failed() error {
	return res.Err
}

// check if AddResponse implements endpoint.Failer in runtime
var _ endpoint.Failer = AddResponse{}

func NewSet(svc svc.AddService, logger log.Logger, duration metrics.Histogram, otTracer stdopentracing.Tracer, zipkinTracer *stdzipkin.Tracer) Set {
	var addEndpoint endpoint.Endpoint
	{
		addEndpoint = MakeSumEndpoint(svc)
		// ratelimit: max 1 request per second
		addEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(addEndpoint)

		addEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(addEndpoint)
		addEndpoint = opentracing.TraceServer(otTracer, "Sum")(addEndpoint)
		if zipkinTracer != nil {
			addEndpoint = zipkin.TraceEndpoint(zipkinTracer, "Sum")(addEndpoint)
		}
		addEndpoint = LoggingMiddleware(logger)(addEndpoint)
		addEndpoint = MetricsMiddleware(duration)(addEndpoint)
	}
	return Set{AddEndpoint: addEndpoint}
}
func (s Set) Add(ctx context.Context, a, b int32) (int32, error) {
	res, err := s.AddEndpoint(ctx, AddRequest{A: a, B: b})
	return res.(AddResponse).V, err
}
func MakeSumEndpoint(svc svc.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(AddRequest)
		v, err := svc.Add(ctx, r.A, r.B)
		return &AddResponse{V: v, Err: err}, err
	}
}
