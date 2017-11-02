package model

var (
	ErrMissingField = "Error missing %v"
)

// ProductCatalog 分类
type ProductCatalog struct {
	ID          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
}

// Product 商品信息
type Product struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Price       string   `json:"price" bson:"price"`
	ID          string   `json:"id" bson:"-"`
	UserID      string   `json:"userID" bson:"userID"`
	CatalogID   string   `json:"catalogID" bson:"catalogID"`
	Status      int      `json:"status" bson:"status"`
	Thumbnails  []string `json:"thumbnails" bson:"thumbnails"`
}

// New a new product instance
func New() Product {
	p := Product{}
	return p
}
