package endpoint

import (
	"context"

	rl "github.com/juju/ratelimit"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	"github.com/laidingqing/dabanshan/svcs/product/service"
	"github.com/laidingqing/dabanshan/svcs/product/model"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	CreateProductEndpoint endpoint.Endpoint
	GetProductsEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, trace stdopentracing.Tracer) Set {
	var (
		createProductEndpoint endpoint.Endpoint
		getProductsEndpoint endpoint.Endpoint 
	)
	{
		createProductEndpoint = MakeCreateProductEndpoint(svc)
		createProductEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(createProductEndpoint)
		createProductEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(createProductEndpoint)
		createProductEndpoint = opentracing.TraceServer(trace, "GetProducts")(createProductEndpoint)
		createProductEndpoint = LoggingMiddleware(log.With(logger, "method", "GetProducts"))(createProductEndpoint)
		createProductEndpoint = InstrumentingMiddleware(duration.With("method", "GetProducts"))(createProductEndpoint)
	}
	{
		getProductsEndpoint = MakeGetProductsEndpoint(svc)
		getProductsEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(getProductsEndpoint)
		getProductsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getProductsEndpoint)
		getProductsEndpoint = opentracing.TraceServer(trace, "GetProducts")(getProductsEndpoint)
		getProductsEndpoint = LoggingMiddleware(log.With(logger, "method", "GetProducts"))(getProductsEndpoint)
		getProductsEndpoint = InstrumentingMiddleware(duration.With("method", "GetProducts"))(getProductsEndpoint)
	}
	return Set{
		GetProductsEndpoint: getProductsEndpoint,
		CreateProductEndpoint: createProductEndpoint,
	}
}

// GetProducts implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) GetProducts(ctx context.Context, a, b int64) (int64, error) {
	resp, err := s.GetProductsEndpoint(ctx, model.GetProductsRequest{A: a, B: b})
	if err != nil {
		return 0, err
	}
	response := resp.(model.GetProductsResponse)
	return response.V, response.Err
}

// CreateProduct implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) CreateProduct(ctx context.Context, req model.CreateProductRequest) (model.CreateProductResponse, error) {
	resp, err := s.CreateProductEndpoint(ctx, req)
	if err != nil {
		return model.CreateProductResponse{}, err
	}
	response := resp.(model.CreateProductResponse)
	return response, response.Err
}


// MakeGetProductsEndpoint constructs a GetProducts endpoint wrapping the service.
func MakeGetProductsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.GetProductsRequest)
		v, err := s.GetProducts(ctx, req.A, req.B)
		return model.GetProductsResponse{V: v, Err: err}, err
	}
}

// MakeCreateProductEndpoint ...
func MakeCreateProductEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(model.CreateProductRequest)
		v, err := s.CreateProduct(ctx, req)
		return v, err
	}
}


