package svc

import "context"

// https://github.com/go-kit/examples/blob/master/addsvc/pkg/addservice/middleware.go
type AddService interface {
	Add(ctx context.Context, a, b int32) (int32, error)
}

type addsvc struct{}

func NewAddService(ctx context.Context) AddService {
	return &addsvc{}
}

func (svc addsvc) Add(ctx context.Context, a, b int32) (int32, error) {
	return a + b, nil
}
