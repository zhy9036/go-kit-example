package svc

import (
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

type Middleware func(AddService) AddService

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(as AddService) AddService {
		return &loggingMiddleware{
			logger: logger,
			next:   as,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   AddService
}

func (mw loggingMiddleware) Add(ctx context.Context, a, b int) (rst int, err error) {
	defer func() {
		mw.logger.Log("method", "Add", "a", a, "b", b, "rst", rst, "err", err)
	}()
	return mw.next.Add(ctx, a, b)
}

func MetricsMiddleware(ints metrics.Counter) Middleware {
	return func(as AddService) AddService {
		return &metricsMiddleware{
			ints: ints,
			next: as,
		}
	}
}

type metricsMiddleware struct {
	ints metrics.Counter
	next AddService
}

func (mw metricsMiddleware) Add(ctx context.Context, a, b int) (int, error) {
	v, e := mw.next.Add(ctx, a, b)
	mw.ints.Add(float64(v))
	return v, e
}
