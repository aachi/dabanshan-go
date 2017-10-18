package products

import (
	"io"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/laidingqing/dabanshan/proto/product"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/robjsliwa/stringsvc1"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var productCli product.ProductRpcServiceClient

func Init(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) {
	productCli = NewProductClient(conn, tracer, logger)
}

func InitWithSD(sdClient consul.Client, tracer stdopentracing.Tracer, logger log.Logger) {
	productCli = NewProductClientWithSD(sdClient, tracer, logger)
}

func GetClient() product.ProductRpcServiceClient {
	if productCli == nil {
		panic("product client is not be initialized!")
	}
	return productCli
}

type ProductClient struct {
	GetProductsEndpoint endpoint.Endpoint
}

func (f *ProductClient) GetProducts(ctx context.Context, in *product.GetProductsRequest, opts ...grpc.CallOption) (*product.GetProductsResponse, error) {
	resp, err := f.GetProductsEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*product.GetProductsResponse), nil
}

func NewProductClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) product.ProductRpcServiceClient {

	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var getProductsEndpoint endpoint.Endpoint
	{
		getProductsEndpoint = grpctransport.NewClient(
			conn,
			"product.Product",
			"GetProducts",
			stringsvc1.EncodeGRPCUppercaseRequest,
			stringsvc1.DecodeGRPCUppercaseResponse,
			product.GetProductsResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()

		getProductsEndpoint = opentracing.TraceClient(tracer, "GetProducts")(getProductsEndpoint)
		getProductsEndpoint = limiter(getProductsEndpoint)
		getProductsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetProducts",
			Timeout: 5 * time.Second,
		}))(getProductsEndpoint)
	}

	return &ProductClient{
		GetProductsEndpoint: getProductsEndpoint,
	}
}

func MakeGetProductsEndpoint(f product.ProductRpcServiceClient) endpoint.Endpoint {
	return f.(*ProductClient).GetProductsEndpoint
}

func NewProductClientWithSD(sdClient consul.Client, tracer stdopentracing.Tracer, logger log.Logger) product.ProductRpcServiceClient {
	res := &ProductClient{}
	var (
		consulService = "productService"
		consulTags    = []string{"prod"}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)
	factory := ProductFactory(MakeGetProductsEndpoint, tracer, logger)
	instancer := consul.NewInstancer(sdClient, logger, consulService, consulTags, passingOnly)
	endpointer := sd.NewEndpointer(instancer, factory, logger)
	balancer := lb.NewRoundRobin(endpointer)
	retry := lb.Retry(retryMax, retryTimeout, balancer)
	res.GetProductsEndpoint = retry

	return res
}

func ProductFactory(makeEndpoint func(f product.ProductRpcServiceClient) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := NewProductClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)

		return endpoint, conn, nil
	}
}
