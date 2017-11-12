package transport

import (
	"net/http"
	"errors"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	// "github.com/go-kit/kit/examples/addsvc/pkg/addendpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	// "github.com/laidingqing/dabanshan/svcs/product/service"
)

var (
	// ErrBadRouting ..
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints p_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	// m := http.NewServeMux()
	r := mux.NewRouter()

	createProductHandle := httptransport.NewServer(
		endpoints.CreateProductEndpoint,
		decodeHTTPCreateProductRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateProduct", logger)))...,
	)

	listProductHandle := httptransport.NewServer(
		endpoints.GetProductsEndpoint,
		decodeHTTPGetProductRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetProducts", logger)))...,
	)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Log("params", r.FormValue("user"))
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/api/v1/products/", listProductHandle).Methods("GET") //获取所有商品，包含按条件分页:catalogID=?
	r.Handle("/api/v1/products/{id}", nil).Methods("GET")           //根据ID获取指定商品
	r.Handle("/api/v1/products/{id}", nil).Methods("DELETE")        //下架指定商品
	r.Handle("/api/v1/products/{id}", nil).Methods("PUT")           //修改指定商品
	r.Handle("/api/v1/products/create", createProductHandle).Methods("POST")           //新增商品
	return r
}






