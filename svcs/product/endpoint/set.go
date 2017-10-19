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
	"github.com/laidingqing/dabanshan/pb"
	"github.com/laidingqing/dabanshan/svcs/product/service"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	GetProductsEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, trace stdopentracing.Tracer) Set {
	var getProductsEndpoint endpoint.Endpoint
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
	}
}

// GetProducts implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) GetProducts(ctx context.Context, req *pb.GetProductsRequest) (pb.GetProductsResponse, error) {
	resp, err := s.GetProductsEndpoint(ctx, req)
	response := resp.(pb.GetProductsResponse)
	return response, err
}

// MakeGetProductsEndpoint constructs a GetProducts endpoint wrapping the service.
func MakeGetProductsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetProductsRequest)
		v, err := s.GetProducts(ctx, &pb.GetProductsRequest{
			Creatorid: req.A,
			Size:      req.B,
		})
		return GetProductsResponse{V: v.GetProducts(), Err: err}, nil
	}
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// GetProductsRequest collects the request parameters for the GetProducts method.
type GetProductsRequest struct {
	A, B int64
}

// GetProductsResponse collects the response values for the GetProducts method.
type GetProductsResponse struct {
	V   []*pb.ProductRecord `json:"v"`
	Err error               `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r GetProductsResponse) Failed() error { return r.Err }
