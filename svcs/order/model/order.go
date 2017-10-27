package model

// OrderItem represents .
type OrderItem struct {
	Quantity  int     `json:"quantity" bson:"quantity"`
	ProductID int64   `json:"code" bson:"productId"`
	Price     float32 `json:"price" bson:"price"`
	Total     float32 `json:"total" bson:"total"`
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

// New ..
func New() Invoice {
	u := Invoice{}
	return u
}

// CreateOrderRequest struct
type CreateOrderRequest struct {
	Amount float32 `json:"amount"`
}

// CreatedOrderResponse ...
type CreatedOrderResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

type Failer interface {
	Failed() error
}
