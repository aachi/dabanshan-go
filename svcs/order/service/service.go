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
	GetOrders(ctx context.Context, req model.GetOrdersRequest) (model.GetOrdersResponse, error)
	AddCart(ctx context.Context, req model.CreateCartRequest) (model.CreatedCartResponse, error)
	GetCartItems(ctx context.Context, req model.GetCartItemsRequest) (model.GetCartItemsResponse, error)
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
func (s basicService) CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreatedOrderResponse, error) {
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

// GetOrders get orders by user id
func (s basicService) GetOrders(ctx context.Context, req model.GetOrdersRequest) (model.GetOrdersResponse, error) {
	orders, err := db.GetOrders(req.UserID)
	if err != nil {
		return model.GetOrdersResponse{Err: err}, err
	}
	return model.GetOrdersResponse{
		Orders: orders,
		Err:    nil,
	}, nil
}

// GetUser get user by id
func (s basicService) AddCart(ctx context.Context, order model.CreateCartRequest) (model.CreatedCartResponse, error) {
	c := model.Cart{}
	c.Price = order.Price
	c.ProductID = order.ProductID
	c.UserID = order.UserID
	// TODO 校验等
	id, err := db.AddCart(&c)
	if err != nil {
		return model.CreatedCartResponse{ID: "", Err: err}, err
	}
	return model.CreatedCartResponse{
		ID:  id,
		Err: nil,
	}, nil
}

// GetCartItems find user's cart items
func (s basicService) GetCartItems(ctx context.Context, req model.GetCartItemsRequest) (model.GetCartItemsResponse, error) {
	items, err := db.GetCartItems(req.UserID)
	if err != nil {
		return model.GetCartItemsResponse{
			Err: err,
		}, err
	}
	return model.GetCartItemsResponse{
		Items: items,
		Err:   nil,
	}, nil
}
