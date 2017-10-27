package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/order/db"
	"github.com/laidingqing/dabanshan/svcs/order/model"
)

var (
	// ErrOrderNotFound ...
	ErrOrderNotFound = errors.New("not found order")
)

// Service describes a service that adds things together.
type Service interface {
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreatedOrderResponse, error)
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
	ErrUnauthorized = errors.New("Unauthorized")
)

const ()

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

// GetUser get user by id
func (s basicService) CreateOrder(_ context.Context, order model.CreateOrderRequest) (model.CreatedOrderResponse, error) {
	u := model.New()
	id, err := db.CreateOrder(&u)
	if err != nil {
		return model.CreatedOrderResponse{ID: "", Err: err}, err
	}
	return model.CreatedOrderResponse{
		ID:  id,
		Err: nil,
	}, nil
}
