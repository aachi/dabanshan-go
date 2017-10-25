package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/examples/addsvc/pkg/addendpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	p_endpoint "github.com/laidingqing/dabanshan/svcs/product/endpoint"
	"github.com/laidingqing/dabanshan/svcs/product/service"
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
	r.Handle("/api/v1/products/", listProductHandle).Methods("GET")
	return r
}

func decodeHTTPGetProductRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// err := json.NewDecoder(r.Body).Decode(&req)
	// todo convert params..
	a, _ := strconv.ParseInt(r.FormValue("userid"), 10, 64)
	b, _ := strconv.ParseInt(r.FormValue("size"), 10, 64)
	return p_endpoint.GetProductsRequest{A: a, B: b}, nil
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(addendpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func err2code(err error) int {
	switch err {
	case service.ErrTwoZeroes, service.ErrMaxSizeExceeded, service.ErrIntOverflow:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
