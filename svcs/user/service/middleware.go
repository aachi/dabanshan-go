package service

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/user/model"
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

func (mw loggingMiddleware) GetUser(ctx context.Context, a string) (v model.GetUserResponse, err error) {
	defer func() {
		mw.logger.Log("method", "GetUser", "err", err)
	}()
	return mw.next.GetUser(ctx, a)
}

func (mw loggingMiddleware) Register(ctx context.Context, us model.RegisterRequest) (r model.RegisterUserResponse, err error) {
	defer func() {
		mw.logger.Log("method", "Register", "err", err)
	}()
	return mw.next.Register(ctx, us)
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

func (mw instrumentingMiddleware) GetUser(ctx context.Context, a string) (model.GetUserResponse, error) {
	v, err := mw.next.GetUser(ctx, a)
	return v, err
}

func (mw instrumentingMiddleware) Register(ctx context.Context, us model.RegisterRequest) (r model.RegisterUserResponse, err error) {
	v, err := mw.next.Register(ctx, us)
	return v, err
}
