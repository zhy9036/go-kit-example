package svc

import (
	"context"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/log"
)

// https://github.com/go-kit/examples/blob/master/addsvc/pkg/addservice/middleware.go
type AddService interface {
	Add(ctx context.Context, a, b int32) (int32, error)
}

type addsvc struct{}

func New(logger log.Logger, ints metrics.Counter) AddService {
	var svc AddService
	{
		svc = NewAddService(context.Background())
		svc = LoggingMiddleware(logger)(svc)
		svc = MetricsMiddleware(ints)(svc)
	}
	return svc
}
func NewAddService(ctx context.Context) AddService {
	return &addsvc{}
}

func (svc addsvc) Add(ctx context.Context, a, b int32) (int32, error) {
	return a + b, nil
}
