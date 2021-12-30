package svc

import "context"

type AddService interface {
	Add(ctx context.Context, a, b int) (int, error)
}

type addsvc struct{}

func NewAddService(ctx context.Context) AddService {
	return &addsvc{}
}

func (svc addsvc) Add(ctx context.Context, a, b int) (int, error) {
	return a + b, nil
}
