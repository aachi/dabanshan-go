package transport

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	httptransport "github.com/go-kit/kit/transport/http"
	o_endpoint "github.com/laidingqing/dabanshan/svcs/order/endpoint"
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
	//m := http.NewServeMux()
	r := mux.NewRouter()

	createOrderHandle := httptransport.NewServer(
		endpoints.CreateOrderEndpoint,
		decodeHTTPCreateOrderRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "CreateOrder", logger)))...,
	)

	getOrderHandle := httptransport.NewServer(
		endpoints.GetOrdersEndpoint,
		decodeHTTPGetOrdersRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetOrders", logger)))...,
	)

	addCartHandle := httptransport.NewServer(
		endpoints.CreateCartEndpoint,
		decodeHTTPAddCartRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "AddCart", logger)))...,
	)

	getCartItemsHandle := httptransport.NewServer(
		endpoints.GetCartItemsEndpoint,
		decodeHTTPGetCartItemsRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "GetCartItems", logger)))...,
	)

	removeCartItemHandle := httptransport.NewServer(
		endpoints.RemoveCartItemEndpoint,
		decodeHTTPRemoveCartItemRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "RemoveCartItem", logger)))...,
	)
	updateQuantityHandle := httptransport.NewServer(
		endpoints.UpdateQuantityEndpoint,
		decodeHTTPUpdateQuantityRequest,
		encodeHTTPGenericResponse,
		append(options, httptransport.ServerBefore(opentracing.HTTPToContext(tracer, "UpdateQuantity", logger)))...,
	)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	// r.Handle("/api/v1/orders/", negroni.New(
	// 	negroni.HandlerFunc(authorize.JwtMiddleware.HandlerWithNext),
	// 	negroni.Wrap(createOrderHandle),
	// )).Methods("POST") //创建订单
	r.Handle("/api/v1/orders/", createOrderHandle).Methods("POST") //创建订单
	//r.Handle("/api/v1/orders/{id}/", nil).Methods("POST")                       //更新订单项
	//r.Handle("/api/v1/orders/{id}/", nil).Methods("GET")                        //查看订单详情
	//r.Handle("/api/v1/orders/{id}/", nil).Methods("DELETE")                     //关闭订单
	r.Handle("/api/v1/orders/", getOrderHandle).Methods("GET")                  //查询用户订单订单项 ?userId=xxxx
	r.Handle("/api/v1/carts/", addCartHandle).Methods("POST")                   //添加至购物车
	r.Handle("/api/v1/carts/", getCartItemsHandle).Methods("GET")               //获取所有购物车数据
	r.Handle("/api/v1/carts/{cartId}/", updateQuantityHandle).Methods("PUT")    //更新购物车项数量
	r.Handle("/api/v1/carts/{cartId}/", removeCartItemHandle).Methods("DELETE") //删除购物车内记录
	return r
}
