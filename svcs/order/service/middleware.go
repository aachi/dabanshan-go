package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/order/model"
)

// Middleware ...
type Middleware func(Service) Service

// LoggingMiddleware ..
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{logger, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) CreateOrder(ctx context.Context, a model.CreateOrderRequest) (v model.CreatedOrderResponse, err error) {
	defer func() {
		mw.logger.Log("method", "CreateOrder", "err", err)
	}()
	return mw.next.CreateOrder(ctx, a)
}

func (mw loggingMiddleware) GetOrders(ctx context.Context, a model.GetOrdersRequest) (v model.GetOrdersResponse, err error) {
	defer func() {
		mw.logger.Log("method", "CreateOrder", "err", err)
	}()
	return mw.next.GetOrders(ctx, a)
}

func (mw loggingMiddleware) AddCart(ctx context.Context, a model.CreateCartRequest) (v model.CreatedCartResponse, err error) {
	defer func() {
		mw.logger.Log("method", "AddCart", "err", err)
	}()
	return mw.next.AddCart(ctx, a)
}

func (mw loggingMiddleware) GetCartItems(ctx context.Context, req model.GetCartItemsRequest) (v model.GetCartItemsResponse, err error) {
	defer func() {
		mw.logger.Log("method", "GetCartItems", "userID", req.UserID, "err", err)
	}()
	return mw.next.GetCartItems(ctx, req)
}

func (mw loggingMiddleware) RemoveCartItem(ctx context.Context, req model.RemoveCartItemRequest) (v model.RemoveCartItemResponse, err error) {
	defer func() {
		mw.logger.Log("method", "RemoveCartItem", "cartID", req.CartID, "err", err)
	}()
	return mw.next.RemoveCartItem(ctx, req)
}

// InstrumentingMiddleware ..
func InstrumentingMiddleware(ints, chars metrics.Counter) Middleware {
	return func(next Service) Service {
		return instrumentingMiddleware{
			ints:  ints,
			chars: chars,
			next:  next,
		}
	}
}

type instrumentingMiddleware struct {
	ints  metrics.Counter
	chars metrics.Counter
	next  Service
}

func (mw instrumentingMiddleware) CreateOrder(ctx context.Context, a model.CreateOrderRequest) (model.CreatedOrderResponse, error) {
	v, err := mw.next.CreateOrder(ctx, a)
	return v, err
}

func (mw instrumentingMiddleware) GetOrders(ctx context.Context, a model.GetOrdersRequest) (model.GetOrdersResponse, error) {
	v, err := mw.next.GetOrders(ctx, a)
	return v, err
}

func (mw instrumentingMiddleware) AddCart(ctx context.Context, a model.CreateCartRequest) (model.CreatedCartResponse, error) {
	v, err := mw.next.AddCart(ctx, a)
	return v, err
}
func (mw instrumentingMiddleware) GetCartItems(ctx context.Context, req model.GetCartItemsRequest) (model.GetCartItemsResponse, error) {
	v, err := mw.next.GetCartItems(ctx, req)
	return v, err
}

func (mw instrumentingMiddleware) RemoveCartItem(ctx context.Context, req model.RemoveCartItemRequest) (model.RemoveCartItemResponse, error) {
	v, err := mw.next.RemoveCartItem(ctx, req)
	return v, err
}
