package transport

import (
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/laidingqing/dabanshan/pb"
	o_endpoint "github.com/laidingqing/dabanshan/svcs/order/endpoint"
	"github.com/laidingqing/dabanshan/svcs/order/service"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcServer struct {
	createOrder    grpctransport.Handler
	getOrders      grpctransport.Handler
	getOrder       grpctransport.Handler
	addCart        grpctransport.Handler
	getCartItems   grpctransport.Handler
	removeCartItem grpctransport.Handler
	updateQuantity grpctransport.Handler
}

// NewGRPCServer ...
func NewGRPCServer(endpoints o_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.OrderRpcServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		createOrder: grpctransport.NewServer(
			endpoints.CreateOrderEndpoint,
			decodeGRPCCreateOrderRequest,
			encodeGRPCCreateOrderResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "CreateOrder", logger)))...,
		),
		getOrders: grpctransport.NewServer(
			endpoints.GetOrdersEndpoint,
			decodeGRPCGetOrdersRequest,
			encodeGRPCGetOrdersResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetOrders", logger)))...,
		),
		getOrder: grpctransport.NewServer(
			endpoints.GetOrderEndpoint,
			decodeGRPCGetOrderRequest,
			encodeGRPCGetOrderResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetOrder", logger)))...,
		),
		addCart: grpctransport.NewServer(
			endpoints.CreateCartEndpoint,
			decodeGRPCAddCartRequest,
			encodeGRPCAddCartResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "AddCart", logger)))...,
		),
		getCartItems: grpctransport.NewServer(
			endpoints.GetCartItemsEndpoint,
			decodeGRPCGetCartItemsRequest,
			encodeGRPCGetCartItemsResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetCartItems", logger)))...,
		),
		removeCartItem: grpctransport.NewServer(
			endpoints.RemoveCartItemEndpoint,
			decodeGRPCRemoveCartItemRequest,
			encodeGRPCRemoveCartItemResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "RemoveCartItem", logger)))...,
		),
		updateQuantity: grpctransport.NewServer(
			endpoints.UpdateQuantityEndpoint,
			decodeGRPCUpdateQuantityRequest,
			encodeGRPCUpdateQuantityResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "UpdateQuantity", logger)))...,
		),
	}
}

// GetUser RPC
func (s *grpcServer) CreateOrder(ctx oldcontext.Context, req *pb.CreateOrderRequest) (*pb.CreatedOrderResponse, error) {
	_, rep, err := s.createOrder.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.CreatedOrderResponse)
	return res, nil
}

// GetOrders

func (s *grpcServer) GetOrders(ctx oldcontext.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	_, rep, err := s.getOrders.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetOrdersResponse)
	return res, nil
}

// GetOrder
func (s *grpcServer) GetOrder(ctx oldcontext.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	_, rep, err := s.getOrder.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetOrderResponse)
	return res, nil
}

// AddCart
func (s *grpcServer) AddCart(ctx oldcontext.Context, req *pb.CreateCartRequest) (*pb.CreatedCartResponse, error) {
	_, rep, err := s.addCart.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.CreatedCartResponse)
	return res, nil
}

// GetCartItems
func (s *grpcServer) GetCartItems(ctx oldcontext.Context, req *pb.GetCartItemsRequest) (*pb.GetCartItemsResponse, error) {
	_, rep, err := s.getCartItems.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetCartItemsResponse)
	return res, nil
}

// RemoveCartItem
func (s *grpcServer) RemoveCartItem(ctx oldcontext.Context, req *pb.RemoveCartItemRequest) (*pb.RemoveCartItemResponse, error) {
	_, rep, err := s.removeCartItem.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.RemoveCartItemResponse)
	return res, nil
}

// UpdateQuantity
func (s *grpcServer) UpdateQuantity(ctx oldcontext.Context, req *pb.UpdateQuantityRequest) (*pb.UpdateQuantityResponse, error) {
	_, rep, err := s.updateQuantity.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.UpdateQuantityResponse)
	return res, nil
}

// NewGRPCClient ...
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	var createOrderEndpoint endpoint.Endpoint
	var getOrdersEndpoint endpoint.Endpoint
	var getOrderEndpoint endpoint.Endpoint
	var addCartEndpoint endpoint.Endpoint
	var getCartItemsEndpoint endpoint.Endpoint
	var removeCartItemEndpoint endpoint.Endpoint
	var updateQuantityEndpoint endpoint.Endpoint
	{
		createOrderEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"CreateOrder",
			encodeGRPCCreateOrderRequest,
			decodeGRPCCreateOrderResponse,
			pb.CreatedOrderResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		createOrderEndpoint = opentracing.TraceClient(tracer, "CreateOrder")(createOrderEndpoint)
		createOrderEndpoint = limiter(createOrderEndpoint)
		createOrderEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "CreateOrder",
			Timeout: 30 * time.Second,
		}))(createOrderEndpoint)

		getOrdersEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"GetOrders",
			encodeGRPCGetOrdersRequest,
			decodeGRPCGetOrdersResponse,
			pb.GetOrdersResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getOrdersEndpoint = opentracing.TraceClient(tracer, "GetOrders")(getOrdersEndpoint)
		getOrdersEndpoint = limiter(getOrdersEndpoint)
		getOrdersEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetOrders",
			Timeout: 30 * time.Second,
		}))(getOrdersEndpoint)

		getOrderEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"GetOrder",
			encodeGRPCGetOrderRequest,
			decodeGRPCGetOrderResponse,
			pb.GetOrdersResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getOrderEndpoint = opentracing.TraceClient(tracer, "GetOrders")(getOrderEndpoint)
		getOrderEndpoint = limiter(getOrderEndpoint)
		getOrderEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetOrder",
			Timeout: 30 * time.Second,
		}))(getOrderEndpoint)

		addCartEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"AddCart",
			encodeGRPCAddCartRequest,
			decodeGRPCAddCartResponse,
			pb.CreatedCartResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		addCartEndpoint = opentracing.TraceClient(tracer, "AddCart")(addCartEndpoint)
		addCartEndpoint = limiter(addCartEndpoint)
		addCartEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "AddCart",
			Timeout: 30 * time.Second,
		}))(addCartEndpoint)

		getCartItemsEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"GetCartItems",
			encodeGRPCCartItemsRequest,
			decodeGRPCCartItemsResponse,
			pb.GetCartItemsResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getCartItemsEndpoint = opentracing.TraceClient(tracer, "GetCartItems")(getCartItemsEndpoint)
		getCartItemsEndpoint = limiter(getCartItemsEndpoint)
		getCartItemsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetCartItems",
			Timeout: 30 * time.Second,
		}))(getCartItemsEndpoint)

		removeCartItemEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"RemoveCartItem",
			encodeGRPCRemoveCartItemRequest,
			decodeGRPCRemoveCartItemResponse,
			pb.RemoveCartItemResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		removeCartItemEndpoint = opentracing.TraceClient(tracer, "RemoveCartItem")(removeCartItemEndpoint)
		removeCartItemEndpoint = limiter(removeCartItemEndpoint)
		removeCartItemEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "RemoveCartItem",
			Timeout: 30 * time.Second,
		}))(removeCartItemEndpoint)

		updateQuantityEndpoint = grpctransport.NewClient(
			conn,
			"pb.OrderRpcService",
			"UpdateQuantity",
			encodeGRPCUpdateQuantityRequest,
			decodeGRPCUpdateQuantityResponse,
			pb.UpdateQuantityResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		updateQuantityEndpoint = opentracing.TraceClient(tracer, "UpdateQuantity")(updateQuantityEndpoint)
		updateQuantityEndpoint = limiter(updateQuantityEndpoint)
		updateQuantityEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "UpdateQuantity",
			Timeout: 30 * time.Second,
		}))(updateQuantityEndpoint)
	}
	return o_endpoint.Set{
		CreateOrderEndpoint:    createOrderEndpoint,
		GetOrdersEndpoint:      getOrdersEndpoint,
		GetOrderEndpoint:       getOrderEndpoint,
		CreateCartEndpoint:     addCartEndpoint,
		GetCartItemsEndpoint:   getCartItemsEndpoint,
		RemoveCartItemEndpoint: removeCartItemEndpoint,
		UpdateQuantityEndpoint: updateQuantityEndpoint,
	}
}
