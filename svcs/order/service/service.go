package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/order/db"
	"github.com/laidingqing/dabanshan/svcs/order/model"
	"github.com/laidingqing/dabanshan/utils"
)

var (
	// ErrOrderNotFound ...
	ErrOrderNotFound = errors.New("not found order")
)

// Service describes a service that adds things together.
type Service interface {
	CreateOrder(ctx context.Context, order model.CreateOrderRequest) (model.CreatedOrderResponse, error)
	GetOrders(ctx context.Context, req model.GetOrdersRequest) (model.GetOrdersResponse, error)
	GetOrder(ctx context.Context, req model.GetOrderRequest) (model.GetOrderResponse, error)
	AddCart(ctx context.Context, req model.CreateCartRequest) (model.CreatedCartResponse, error)
	GetCartItems(ctx context.Context, req model.GetCartItemsRequest) (model.GetCartItemsResponse, error)
	RemoveCartItem(ctx context.Context, req model.RemoveCartItemRequest) (model.RemoveCartItemResponse, error)
	UpdateQuantity(ctx context.Context, req model.UpdateQuantityRequest) (model.UpdateQuantityResponse, error)
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
	id, err := db.CreateOrder(&order.Invoice)
	if err != nil {
		return model.CreatedOrderResponse{ID: "", Err: err}, err
	}
	return model.CreatedOrderResponse{
		ID:  id,
		Err: nil,
	}, nil
}

// GetOrders get orders by user id or tenant id
func (s basicService) GetOrders(ctx context.Context, req model.GetOrdersRequest) (model.GetOrdersResponse, error) {

	var orders utils.Pagination
	var err error
	if req.UserID != "" {
		orders, err = db.GetOrdersByUser(req.UserID, utils.Pagination{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
		})
	}

	if err != nil {
		return model.GetOrdersResponse{Err: err}, err
	}

	if req.TenantID != "" {
		orders, err = db.GetOrdersByTenant(req.TenantID, utils.Pagination{
			PageIndex: req.PageIndex,
			PageSize:  req.PageSize,
		})
	}

	return model.GetOrdersResponse{
		UserID:   req.UserID,
		TenantID: req.TenantID,
		Orders:   orders,
		Err:      nil,
	}, nil
}

// GetOrder get order by id
func (s basicService) GetOrder(ctx context.Context, req model.GetOrderRequest) (model.GetOrderResponse, error) {

	order, err := db.GetOrder(req.OrderID)

	if err != nil {
		return model.GetOrderResponse{Err: err}, err
	}

	return model.GetOrderResponse{
		Order: order,
		Err:   nil,
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

// RemoveCartItem remove cart item by id
func (s basicService) RemoveCartItem(ctx context.Context, req model.RemoveCartItemRequest) (model.RemoveCartItemResponse, error) {
	_, err := db.RemoveCartItem(req.CartID)
	if err != nil {
		return model.RemoveCartItemResponse{
			Err: err,
		}, err
	}

	return model.RemoveCartItemResponse{
		Err: nil,
	}, nil
}

func (s basicService) UpdateQuantity(ctx context.Context, req model.UpdateQuantityRequest) (model.UpdateQuantityResponse, error) {
	var cart = model.Cart{
		CartID:   req.CartID,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	cart, err := db.UpdateQuantity(&cart)

	if err != nil {
		return model.UpdateQuantityResponse{
			Err: err,
		}, err
	}

	return model.UpdateQuantityResponse{}, nil
}
