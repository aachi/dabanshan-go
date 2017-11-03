package model

// OrderItem represents .
type OrderItem struct {
	Quantity  int     `json:"quantity" bson:"quantity"`
	ProductID int64   `json:"code" bson:"productId"`
	Price     float32 `json:"price" bson:"price"`
	Total     float32 `json:"total" bson:"total"`
	CartID    float32 `json:"cartID" bson:"cartID"`
}

// Invoice represents.
type Invoice struct {
	Amount      float32     `json:"amount" bson:"amount"`
	Discount    float32     `json:"discount" bson:"discount"`
	DiscountID  float32     `json:"discountid" bson:"discountId"`
	UserID      string      `json:"userid" bson:"userId"`
	AddressID   string      `json:"addressid" bson:"addressId"`
	OrderedItem []OrderItem `json:"orderedItem" bson:"orderedItem"`
}

// Cart represents.
type Cart struct {
	UserID    string  `json:"userID" bson:"userID"`
	ProductID string  `json:"productID" bson:"productID"`
	Price     float32 `json:"price" bson:"price"`
	CartID    string  `json:"id" bson:"_id"`
}

// New ..
func New() Invoice {
	u := Invoice{}
	return u
}

// CreateOrderRequest struct
type CreateOrderRequest struct {
	Amount float32 `json:"amount"`
}

// CreateCartRequest struct
type CreateCartRequest struct {
	ProductID string  `json:"productID"`
	UserID    string  `json:"userID"`
	Price     float32 `json:"price"`
}

// GetOrdersRequest struct
type GetOrdersRequest struct {
	UserID string `json:"userID"`
}

// CreatedOrderResponse ...
type CreatedOrderResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// CreatedCartResponse ...
type CreatedCartResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// GetOrdersResponse ...
type GetOrdersResponse struct {
	Orders []Invoice `json:"orders"`
	Err    error     `json:"-"`
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

// Failer ...
type Failer interface {
	Failed() error
}
