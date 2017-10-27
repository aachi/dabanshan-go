package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	o_endpoint "github.com/laidingqing/dabanshan/svcs/order/endpoint"
	m_order "github.com/laidingqing/dabanshan/svcs/order/model"
	"github.com/laidingqing/dabanshan/svcs/order/service"
)

var (
	// ErrBadRouting ..
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints o_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorLogger(logger),
	}
	// m := http.NewServeMux()
	r := mux.NewRouter()
	//authenticationMiddleware := authorize.ValidateTokenMiddleware()

	createOrderHandle := httptransport.NewServer(
		endpoints.CreateOrderEndpoint,
		decodeHTTPCreateOrderRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateOrder", logger)))...,
	)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/api/v1/orders/", createOrderHandle).Methods("POST")
	return r
}

func decodeHTTPCreateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	a := m_order.CreateOrderRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
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
	if f, ok := response.(m_order.Failer); ok && f.Failed() != nil {
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
	case service.ErrOrderNotFound:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
