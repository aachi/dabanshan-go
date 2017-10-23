package main

import (
	"flag"
	"fmt"
	corelog "log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/tabwriter"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/laidingqing/dabanshan/svcs/product/db"
	"github.com/laidingqing/dabanshan/svcs/product/db/mongodb"
	lightstep "github.com/lightstep/lightstep-tracer-go"
	"github.com/oklog/oklog/pkg/group"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"

	addpb "github.com/laidingqing/dabanshan/pb"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	p_service "github.com/laidingqing/dabanshan/svcs/product/service"
	p_transport "github.com/laidingqing/dabanshan/svcs/product/transport"
)

func init() {
	db.Register("mongodb", &mongodb.Mongo{})
}

func main() {
	fs := flag.NewFlagSet("productSvc", flag.ExitOnError)
	var (
		debugAddr      = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr       = fs.String("http-addr", ":8081", "HTTP listen address")
		grpcAddr       = fs.String("grpc-addr", ":8082", "gRPC listen address")
		consulAddr     = flag.String("consul.addr", "localhost:8500", "Consul agent address")
		zipkinURL      = fs.String("zipkin-url", "http://localhost:9411/api/v1/spans", "Enable Zipkin tracing via a collector URL e.g. http://localhost:9411/api/v1/spans")
		lightstepToken = flag.String("lightstep-token", "", "Enable LightStep tracing via a LightStep access token")
		appdashAddr    = flag.String("appdash-addr", "", "Enable Appdash tracing via an Appdash server host:port")
		serviceName    = flag.String("service.name", "productsvc", "Name of the service")
		instance       = flag.Int("instance", 1, "The instance count of the status service")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	fs.Parse(os.Args[1:])

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	controlSvc := &service{
		HTTPAddress: httpAddr,
		GRPCAddress: grpcAddr,
		Name:        serviceName,
		Instance:    *instance}

	var (
		kitconsul consulsd.Client
	)
	{
		var err error
		kitconsul, err = createConsulClient(consulAddr, logger)
		if err != nil {
			logger.Log("err", err)
		}
		err = registerService(kitconsul, controlSvc)
		if err != nil {
			logger.Log("err", err)
		}
	}
	// Determine which tracer to use. We'll pass the tracer to all the
	// components that use it, as a dependency.
	var tracer stdopentracing.Tracer
	{
		if *zipkinURL != "" {
			logger.Log("tracer", "Zipkin", "URL", *zipkinURL)
			collector, err := zipkin.NewHTTPCollector(*zipkinURL)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
			defer collector.Close()
			var (
				debug       = false
				hostPort    = "localhost:80"
				serviceName = "productsvc"
			)
			recorder := zipkin.NewRecorder(collector, debug, hostPort, serviceName)
			tracer, err = zipkin.NewTracer(recorder)
			if err != nil {
				logger.Log("err", err)
				os.Exit(1)
			}
		} else if *lightstepToken != "" {
			logger.Log("tracer", "LightStep") // probably don't want to print out the token :)
			tracer = lightstep.NewTracer(lightstep.Options{
				AccessToken: *lightstepToken,
			})
			defer lightstep.FlushLightStepTracer(tracer)
		} else if *appdashAddr != "" {
			logger.Log("tracer", "Appdash", "addr", *appdashAddr)
			tracer = appdashot.NewTracer(appdash.NewRemoteCollector(*appdashAddr))
		} else {
			logger.Log("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	dbconn := false
	for !dbconn {
		err := db.Init()
		if err != nil {
			if err == db.ErrNoDatabaseSelected {
				corelog.Fatal(err)
			}
			corelog.Print(err)
		} else {
			dbconn = true
		}
	}

	// Create the (sparse) metrics we'll use in the service. They, too, are
	// dependencies that we pass to components that use them.
	var ints, chars metrics.Counter
	{
		// Business-level metrics.
		ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "dabanshan",
			Subsystem: "products",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})
		chars = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "dabanshan",
			Subsystem: "products",
			Name:      "characters_concatenated",
			Help:      "Total count of characters concatenated via the Concat method.",
		}, []string{})
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "dabanshan",
			Subsystem: "products",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())
	http.DefaultServeMux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	var (
		service     = p_service.New(logger, ints, chars)
		endpoints   = p_endpoint.New(service, logger, duration, tracer)
		httpHandler = p_transport.NewHTTPHandler(endpoints, tracer, logger)
		grpcServer  = p_transport.NewGRPCServer(endpoints, tracer, logger)
	)

	var g group.Group
	{
		debugListener, err := net.Listen("tcp", *debugAddr)
		if err != nil {
			logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// The gRPC listener mounts the Go kit gRPC server we created.
		grpcListener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			logger.Log("transport", "gRPC", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "gRPC", "addr", *grpcAddr)
			baseServer := grpc.NewServer()
			addpb.RegisterProductRpcServiceServer(baseServer, grpcServer)
			return baseServer.Serve(grpcListener)
		}, func(error) {
			grpcListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

type service struct {
	GRPCAddress *string
	HTTPAddress *string
	Instance    int
	Name        *string
}

func createConsulClient(consulAddr *string, logger log.Logger) (consulsd.Client, error) {
	consulConfig := api.DefaultConfig()
	if len(*consulAddr) > 0 {
		consulConfig.Address = *consulAddr
	}
	consulClient, err := api.NewClient(consulConfig)
	return consulsd.NewClient(consulClient), err
}

func registerService(client consulsd.Client, svc *service) error {
	check := &api.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://127.0.0.1%v/health", *svc.HTTPAddress),
		Interval: "10s",
		Timeout:  "3s",
	}
	host, strPort, _ := net.SplitHostPort(*svc.GRPCAddress)
	port, _ := strconv.Atoi(strPort)
	reg := &api.AgentServiceRegistration{
		Name:    *svc.Name,
		Address: host,
		Port:    port,
		ID:      *svc.Name + "-" + strconv.Itoa(svc.Instance),
		Tags:    []string{"grpc"},
		Check:   check,
	}
	err := client.Register(reg)
	return err
}
