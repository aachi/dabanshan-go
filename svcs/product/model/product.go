package model

var (
	ErrMissingField = "Error missing %v"
)

type Product struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Price       string `json:"price" bson:"price"`
	ProductID   string `json:"id" bson:"-"`
	UserID      string `json:"userid" bson:"userid"`
}

func New() Product {
	p := Product{}
	return p
}
