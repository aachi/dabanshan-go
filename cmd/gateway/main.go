package main

import (
	"flag"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	p_service "github.com/laidingqing/dabanshan/svcs/product/service"
	p_transport "github.com/laidingqing/dabanshan/svcs/product/transport"

	u_endpoint "github.com/laidingqing/dabanshan/svcs/user/endpoint"
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
	)
	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

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
	// ctx := context.Background()
	// r := mux.NewRouter()
	mux := http.NewServeMux()
	// products routes.
	{
		var (
			tags        = []string{}
			passingOnly = true
			endpoints   = p_endpoint.Set{}
			instancer   = consulsd.NewInstancer(client, logger, "productsvc", tags, passingOnly)
		)
		{
			factory := addsvcFactory(p_endpoint.MakeGetProductsEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(instancer, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			endpoints.GetProductsEndpoint = retry
		}
		{
			factory := addsvcFactory(u_endpoint.MakeGetUserEndpoint, tracer, logger)
			endpointer := sd.NewEndpointer(instancer, factory, logger)
			balancer := lb.NewRoundRobin(endpointer)
			retry := lb.Retry(*retryMax, *retryTimeout, balancer)
			endpoints.GetUserEndpoint = retry
		}
		mux.Handle("/api/products", p_transport.NewHTTPHandler(endpoints, tracer, logger))
		mux.HandleFunc("/api/echo", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.FormValue("user")))
		})
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

func addsvcFactory(makeEndpoint func(p_service.Service) endpoint.Endpoint, tracer stdopentracing.Tracer, logger log.Logger) sd.Factory {
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
