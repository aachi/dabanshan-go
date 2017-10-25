package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/user/model"
)

// Service describes a service that adds things together.
type Service interface {
	GetUser(ctx context.Context, id string) (model.GetUserResponse, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger, ints, chars metrics.Counter) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

var (
	// ErrTwoZeroes ..
	ErrTwoZeroes = errors.New("can't sum two zeroes")
	// ErrIntOverflow ...
	ErrIntOverflow = errors.New("integer overflow")
	// ErrMaxSizeExceeded ...
	ErrMaxSizeExceeded = errors.New("result exceeds maximum size")
)

const (
	intMax = 1<<31 - 1
	intMin = -(intMax + 1)
	maxLen = 10
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (s basicService) GetUser(_ context.Context, id string) (model.GetUserResponse, error) {
	return model.GetUserResponse{
		V:   model.New(),
		Err: nil,
	}, nil
}
