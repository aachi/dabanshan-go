package gateway

import (
	"flag"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/laidingqing/dabanshan/client/products"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	var (
		httpAddr   = flag.String("http.addr", ":8080", "HTTP server address")
		consulAddr = flag.String("consul.addr", "", "consul registry address")
		zipkinAddr = flag.String("zipkin.addr", "", "tracer server address")
	)
	flag.Parse()
	// Logging domain.
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	var sdClient consul.Client

	apiclient, err := consulapi.NewClient(&consulapi.Config{
		Address: *consulAddr,
	})

	sdClient = consul.NewClient(apiclient)

	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}

	// Transport domain.
	tracer := stdopentracing.GlobalTracer() // nop by default
	if *zipkinAddr != "" {
		logger := log.With(logger, "tracer", "Zipkin")
		logger.Log("addr", *zipkinAddr)
		collector, err := zipkin.NewHTTPCollector(
			*zipkinAddr,
			zipkin.HTTPLogger(logger),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		tracer, err = zipkin.NewTracer(
			zipkin.NewRecorder(collector, false, "localhost:80", "http"),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
	}

	// Debug listener.
	go func() {
		logger := log.With(logger, "transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", stdprometheus.Handler())

		logger.Log("addr", ":6060")
		http.ListenAndServe(":6060", m)
	}()

	products.InitWithSD(sdClient, tracer, logger)
	//Inject
	router := gin.New()
	Register(router)

	server := &http.Server{Addr: *httpAddr, Handler: router}
	if err = gracehttp.Serve(server); err != nil {
		panic(err)
	}
}
