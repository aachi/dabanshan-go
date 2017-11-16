package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/laidingqing/dabanshan/svcs/order/model"
	"github.com/laidingqing/dabanshan/svcs/order/service"
	"github.com/laidingqing/dabanshan/utils"
)

func decodeHTTPCreateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	logger := utils.NewLogger()

	defer r.Body.Close()
	a := model.CreateOrderRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	logger.Log("amount", a.Invoice.Amount, "userId", a.Invoice.UserID, "items", len(a.Invoice.OrdereItem))
	if err != nil {
		return nil, err
	}
	return a, nil
}

func decodeHTTPGetOrdersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	userID := r.FormValue("userId")
	tenantID := r.FormValue("tenantId")
	pageIndex, _ := strconv.Atoi(r.FormValue("pageIndex"))
	pageSize, _ := strconv.Atoi(r.FormValue("pageSize"))
	logger := utils.NewLogger()
	logger.Log("pageIndex", pageIndex, "pageSize", pageSize, "userID", userID)
	a := model.GetOrdersRequest{
		UserID:    userID,
		TenantID:  tenantID,
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}
	return a, nil
}

func decodeHTTPGetOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := vars["orderId"]
	a := model.GetOrderRequest{
		OrderID: id,
	}
	return a, nil
}

func decodeHTTPAddCartRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	a := model.CreateCartRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func decodeHTTPGetCartItemsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id := r.FormValue("userId")
	return model.GetCartItemsRequest{
		UserID: id,
	}, nil
}

func decodeHTTPRemoveCartItemRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := vars["cartId"]
	return model.RemoveCartItemRequest{
		CartID: id,
	}, nil
}

func decodeHTTPUpdateQuantityRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, _ := vars["cartId"]

	defer r.Body.Close()
	a := model.UpdateQuantityRequest{}
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		return nil, err
	}
	return model.UpdateQuantityRequest{
		CartID:   id,
		Quantity: a.Quantity,
	}, nil
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
	if f, ok := response.(model.Failer); ok && f.Failed() != nil {
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
