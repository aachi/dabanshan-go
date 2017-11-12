package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/product/model"
)

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

func (mw loggingMiddleware) GetProducts(ctx context.Context, a, b int64) (v int64, err error) {
	defer func() {
		mw.logger.Log("method", "GetProducts", "err", err)
	}()
	return mw.next.GetProducts(ctx, a, b)
}

func (mw loggingMiddleware) CreateProduct(ctx context.Context, req model.CreateProductRequest) (res model.CreateProductResponse, err error) {
	defer func() {
		mw.logger.Log("method", "CreateProduct", "err", err)
	}()
	return mw.next.CreateProduct(ctx, req)
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

func (mw instrumentingMiddleware) GetProducts(ctx context.Context, a, b int64) (int64, error) {
	v, err := mw.next.GetProducts(ctx, a, b)
	// mw.ints.Add(float64(v))
	return v, err
}

func (mw instrumentingMiddleware) CreateProduct(ctx context.Context, req model.CreateProductRequest) (model.CreateProductResponse, error) {
	v, err := mw.next.CreateProduct(ctx, req)
	return v, err
}

