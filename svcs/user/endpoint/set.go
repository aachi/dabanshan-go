package endpoint

import (
	"context"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	rl "github.com/juju/ratelimit"
	m_user "github.com/laidingqing/dabanshan/svcs/user/model"
	"github.com/laidingqing/dabanshan/svcs/user/service"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
)

// Set collects all of the endpoints that compose an add service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Set struct {
	GetUserEndpoint endpoint.Endpoint
}

// New returns a Set that wraps the provided server, and wires in all of the
// expected endpoint middlewares via the various parameters.
func New(svc service.Service, logger log.Logger, duration metrics.Histogram, trace stdopentracing.Tracer) Set {
	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = MakeGetUserEndpoint(svc)
		getUserEndpoint = ratelimit.NewTokenBucketLimiter(rl.NewBucketWithRate(1, 1))(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getUserEndpoint)
		getUserEndpoint = opentracing.TraceServer(trace, "GetUser")(getUserEndpoint)
		getUserEndpoint = LoggingMiddleware(log.With(logger, "method", "GetUser"))(getUserEndpoint)
		getUserEndpoint = InstrumentingMiddleware(duration.With("method", "GetUser"))(getUserEndpoint)
	}

	return Set{
		GetUserEndpoint: getUserEndpoint,
	}
}

// GetUser implements the service interface, so Set may be used as a service.
// This is primarily useful in the context of a client library.
func (s Set) GetUser(ctx context.Context, a string) (m_user.GetUserResponse, error) {
	resp, err := s.GetUserEndpoint(ctx, m_user.GetUserRequest{A: a})
	if err != nil {
		return m_user.GetUserResponse{}, err
	}
	response := resp.(m_user.GetUserResponse)
	return response, response.Err
}

// MakeGetUserEndpoint constructs a GetProducts endpoint wrapping the service.
func MakeGetUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(m_user.GetUserRequest)
		v, err := s.GetUser(ctx, req.A)
		return v, err
	}
}
