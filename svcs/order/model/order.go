package model

import (
	"time"

	"github.com/laidingqing/dabanshan/utils"
)

// OrderItem represents .
type OrderItem struct {
	Quantity  int32   `json:"quantity" bson:"quantity"`
	ProductID string  `json:"code" bson:"productId"`
	Price     float32 `json:"price" bson:"price"`
	Total     float32 `json:"total" bson:"total"`
	CartID    string  `json:"cartID" bson:"cartID"`
	TenantID  string  `json:"tenantId" bson:"tenantId"`
}

// Invoice represents.
type Invoice struct {
	Amount     float32     `json:"amount" bson:"amount"`
	Discount   float32     `json:"discount" bson:"discount"`
	DiscountID float32     `json:"discountid" bson:"discountId"`
	UserID     string      `json:"userid" bson:"userId"`
	AddressID  string      `json:"addressId" bson:"addressId"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	Status     OrderStatus `json:"status" bson:"status"`
	TenantID   string      `json:"tenantID" bson:"tenantID"`
	OrdereItem []OrderItem `json:"items" bson:"items"`
}

// Procurement represents. 采购清单
type Procurement struct {
	Amount     float32     `json:"amount" bson:"amount"`
	UserID     string      `json:"userid" bson:"userId"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt" bson:"updatedAt"`
	OrdereItem []OrderItem `json:"items" bson:"items"`
}

// Cart represents.
type Cart struct {
	UserID    string  `json:"userID" bson:"userID"`
	ProductID string  `json:"productID" bson:"productID"`
	Price     float32 `json:"price" bson:"price"`
	Quantity  int32   `json:"quantity" bson:"quantity"`
	CartID    string  `json:"id" bson:"-"`
	Total     float32 `json:"total" bson:"total"`
}

// New ..
func New() Invoice {
	u := Invoice{}
	return u
}

// CreateOrderRequest struct
type CreateOrderRequest struct {
	Invoice Invoice `json:"invoice"`
}

// CreatedOrderResponse ...
type CreatedOrderResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// CreateCartRequest struct
type CreateCartRequest struct {
	ProductID string  `json:"productID"`
	UserID    string  `json:"userID"`
	Price     float32 `json:"price"`
}

// GetOrdersRequest struct
type GetOrdersRequest struct {
	UserID    string `json:"userID"`
	TenantID  string `json:"TenantID"`
	PageIndex int    `json:"pageIndex"`
	PageSize  int    `json:"pageSize"`
}

// GetOrderRequest struct
type GetOrderRequest struct {
	OrderID string `json:"orderID"`
}

// CreatedCartResponse ...
type CreatedCartResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// GetOrdersResponse ...
type GetOrdersResponse struct {
	UserID   string           `json:"userId"`
	TenantID string           `json:"tenanId"`
	Orders   utils.Pagination `json:"orders"`
	Err      error            `json:"-"`
}

// GetOrderResponse ...
type GetOrderResponse struct {
	Order Invoice `json:"order"`
	Err   error   `json:"-"`
}

// GetCartItemsRequest ...
type GetCartItemsRequest struct {
	UserID string `json:"userID"`
}

// GetCartItemsResponse ..
type GetCartItemsResponse struct {
	Items []Cart `json:"items"`
	Err   error  `json:"-"`
}

// RemoveCartItemRequest ..
type RemoveCartItemRequest struct {
	CartID string
}

// RemoveCartItemResponse ..
type RemoveCartItemResponse struct {
	Err error `json:"-"`
}

// UpdateCartItemRequest ..
type UpdateCartItemRequest struct {
	CartID   string  `json:"cartID"`
	Quantity int32   `json:"quantity"`
	Price    float32 `json:"price"`
}

// UpdateCartItemResponse ..
type UpdateCartItemResponse struct {
	Err error `json:"-"`
}

// UpdateQuantityRequest ...
type UpdateQuantityRequest struct {
	CartID   string  `json:"cartID"`
	Quantity int32   `json:"quantity"`
	Price    float32 `json:"price"`
}

// UpdateQuantityResponse ...
type UpdateQuantityResponse struct {
	Err error `json:"-"`
}

// Failer ...
type Failer interface {
	Failed() error
}
