package transport

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/laidingqing/dabanshan/pb"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	"github.com/laidingqing/dabanshan/svcs/product/service"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	getproducts grpctransport.Handler
}

// NewGRPCServer ...
func NewGRPCServer(endpoints p_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.ProductRpcServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		getproducts: grpctransport.NewServer(
			endpoints.GetProductsEndpoint,
			decodeGRPCGetProductsRequest,
			encodeGRPCGetProductsResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetProducts", logger)))...,
		),
	}
}

func (s *grpcServer) GetProducts(ctx oldcontext.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	_, rep, err := s.getproducts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetProductsResponse)
	return res, nil
}

func decodeGRPCGetProductsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetProductsRequest)
	return p_endpoint.GetProductsRequest{A: int64(req.Creatorid), B: int64(req.Size)}, nil
}

func encodeGRPCGetProductsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(p_endpoint.GetProductsResponse)
	return &pb.GetProductsResponse{V: int64(resp.V), Err: err2str(resp.Err)}, nil
}

// NewGRPCClient ...
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	var getProductsEndpoint endpoint.Endpoint
	{
		getProductsEndpoint = grpctransport.NewClient(
			conn,
			"pb.ProductRpcService",
			"GetProducts",
			encodeGRPCGetProductsRequest,
			decodeGRPCGetProductsResponse,
			pb.GetProductsResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getProductsEndpoint = opentracing.TraceClient(tracer, "GetProducts")(getProductsEndpoint)
		getProductsEndpoint = limiter(getProductsEndpoint)
		getProductsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetProducts",
			Timeout: 30 * time.Second,
		}))(getProductsEndpoint)
	}
	return p_endpoint.Set{
		GetProductsEndpoint: getProductsEndpoint,
	}
}

func encodeGRPCGetProductsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(p_endpoint.GetProductsRequest)
	return &pb.GetProductsRequest{Creatorid: int64(req.A), Size: int64(req.B)}, nil
}

func decodeGRPCGetProductsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetProductsResponse)
	return p_endpoint.GetProductsResponse{V: int64(reply.V), Err: str2err(reply.Err)}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
