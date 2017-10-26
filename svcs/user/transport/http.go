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
	p_endpoint "github.com/laidingqing/dabanshan/svcs/user/endpoint"
	m_user "github.com/laidingqing/dabanshan/svcs/user/model"
	"github.com/laidingqing/dabanshan/svcs/user/service"
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

	getUserHandle := httptransport.NewServer(
		endpoints.GetUserEndpoint,
		decodeHTTPGetUserRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetUser", logger)))...,
	)

	registerHandle := httptransport.NewServer(
		endpoints.GetUserEndpoint,
		decodeHTTPRegisterRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetUser", logger)))...,
	)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Handle("/api/v1/users/{id}", getUserHandle).Methods("GET")
	r.Handle("/api/v1/users/", registerHandle).Methods("POST")
	return r
}

func decodeHTTPRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	a := m_user.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func decodeHTTPGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return m_user.GetUserRequest{A: id}, nil
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
	if f, ok := response.(m_user.Failer); ok && f.Failed() != nil {
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
	case service.ErrUserNotFound:
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
