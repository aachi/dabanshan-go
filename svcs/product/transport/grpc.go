package transport

import (
	"time"
	"fmt"

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
	createProduct grpctransport.Handler
	getproducts grpctransport.Handler
}

// NewGRPCServer ...
func NewGRPCServer(endpoints p_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.ProductRpcServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		createProduct: grpctransport.NewServer(
			endpoints.CreateProductEndpoint,
			decodeGRPCCreateProductRequest,
			encodeGRPCCreateProductResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "CreateProduct", logger)))...,
		),
		getproducts: grpctransport.NewServer(
			endpoints.GetProductsEndpoint,
			decodeGRPCGetProductsRequest,
			encodeGRPCGetProductsResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetProducts", logger)))...,
		),
	}
}

// get products
func (s *grpcServer) GetProducts(ctx oldcontext.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	_, rep, err := s.getproducts.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetProductsResponse)
	return res, nil
}

// create product
func (s *grpcServer) CreateProduct(ctx oldcontext.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	fmt.Println("create name fmt")
	_, rep, err := s.createProduct.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.CreateProductResponse)
	return res, nil
}

// NewGRPCClient ...
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	var getProductsEndpoint endpoint.Endpoint
	var createProductEndpoint endpoint.Endpoint
	{
		createProductEndpoint = grpctransport.NewClient(
			conn,
			"pb.ProductRpcService",
			"CreateProduct",
			encodeGRPCCreateProductRequest,
			decodeGRPCCreateProductResponse,
			pb.CreateProductResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		createProductEndpoint = opentracing.TraceClient(tracer, "CreateProduct")(createProductEndpoint)
		createProductEndpoint = limiter(createProductEndpoint)
		createProductEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateProduct",
			Timeout: 30 * time.Second,
		}))(createProductEndpoint)
	}
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
		CreateProductEndpoint: createProductEndpoint,
		GetProductsEndpoint: getProductsEndpoint,
	}
}




