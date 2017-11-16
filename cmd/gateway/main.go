package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	p_service "github.com/laidingqing/dabanshan/svcs/product/service"
	p_transport "github.com/laidingqing/dabanshan/svcs/product/transport"
	"github.com/laidingqing/dabanshan/utils"
	"google.golang.org/grpc"

	u_endpoint "github.com/laidingqing/dabanshan/svcs/user/endpoint"
	u_service "github.com/laidingqing/dabanshan/svcs/user/service"
	u_transport "github.com/laidingqing/dabanshan/svcs/user/transport"

	o_endpoint "github.com/laidingqing/dabanshan/svcs/order/endpoint"
	o_service "github.com/laidingqing/dabanshan/svcs/order/service"
	o_transport "github.com/laidingqing/dabanshan/svcs/order/transport"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
)

func main() {
	var (
		httpAddr     = flag.String("http.addr", ":8000", "Address for HTTP (JSON) server")
		consulAddr   = flag.String("consul.addr", "localhost:8500", "Consul agent address")
		retryMax     = flag.Int("retry.max", 3, "per-request retries to different instances")
		retryTimeout = flag.Duration("retry.timeout", 500*time.Millisecond, "per-request timeout, including retries")
		staticDir    = flag.String("static_dir", "./public/", "static directory in addition to default static directory")
	)
	flag.Parse()

	// Logging domain.
	logger := utils.NewLogger()

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if len(*consulAddr) > 0 {
			consulConfig.Address = *consulAddr
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	// Transport domain.
	tracer := stdopentracing.GlobalTracer() // no-op
	mux := http.NewServeMux()
	//r := mux.NewRouter()
	{
		var (
			tags             = []string{}
			passingOnly      = true
			pEndpoints       = p_endpoint.Set{}
			uEndpoints       = u_endpoint.Set{}
			oEndpoints       = o_endpoint.Set{}
			productInstancer = consulsd.NewInstancer(client, logger, "productsvc", tags, passingOnly)
			userInstancer    = consulsd.NewInstancer(client, logger, "usersvc", tags, passingOnly)
			orderInstancer   = consulsd.NewInstancer(client, logger, "ordersvc", tags, passingOnly)
		)
		{
			productfactory := addProductFactory(p_endpoint.MakeGetProductsEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(productInstancer, productfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			pEndpoints.GetProductsEndpoint = retry
		}
		{
			productfactory := addProductFactory(p_endpoint.MakeCreateProductEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(productInstancer, productfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			pEndpoints.CreateProductEndpoint = retry
		}
		{
			userfactory := addUserFactory(u_endpoint.MakeGetUserEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(userInstancer, userfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			uEndpoints.GetUserEndpoint = retry
		}
		{
			userfactory := addUserFactory(u_endpoint.MakeRegisterEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(userInstancer, userfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			uEndpoints.RegisterEndpoint = retry
		}
		{
			userfactory := addUserFactory(u_endpoint.MakeLoginEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(userInstancer, userfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			uEndpoints.LoginEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeAddCartEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.CreateCartEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeCreateOrderEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.CreateOrderEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeGetCartItemsEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.GetCartItemsEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeRemoveCartItemEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.RemoveCartItemEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeUpdateQuantityEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.UpdateQuantityEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeGetOrdersEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.GetOrdersEndpoint = retry
		}
		{
			orderfactory := addOrderFactory(o_endpoint.MakeGetOrderEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(orderInstancer, orderfactory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			oEndpoints.GetOrderEndpoint = retry
		}

		mux.Handle("/api/v1/products/", p_transport.NewHTTPHandler(pEndpoints, tracer, logger))
		mux.Handle("/api/v1/users/", u_transport.NewHTTPHandler(uEndpoints, tracer, logger))
		mux.Handle("/api/v1/orders/", o_transport.NewHTTPHandler(oEndpoints, tracer, logger))
		mux.Handle("/api/v1/carts/", o_transport.NewHTTPHandler(oEndpoints, tracer, logger))
		mux.Handle("/", http.FileServer(http.Dir(*staticDir)))
	}
	http.Handle("/", accessControl(mux))
	// Interrupt handler.
	errc := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// HTTP transport.
	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errc <- http.ListenAndServe(*httpAddr, nil)
	}()

	// Run!
	logger.Log("exit", <-errc)

}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func addProductFactory(makeEndpoint func(p_service.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := p_transport.NewGRPCClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)
		return endpoint, conn, nil
	}
}

func addOrderFactory(makeEndpoint func(o_service.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := o_transport.NewGRPCClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)
		return endpoint, conn, nil
	}
}

func addUserFactory(makeEndpoint func(u_service.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		conn, err := grpc.Dial(instance, grpc.WithInsecure())
		if err != nil {
			return nil, nil, err
		}
		service := u_transport.NewGRPCClient(conn, tracer, logger)
		endpoint := makeEndpoint(service)
		return endpoint, conn, nil
	}
}
