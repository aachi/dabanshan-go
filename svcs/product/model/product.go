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
	TenantID    string   `json:"tenantID" bson:"tenantID"`
	CatalogID   string   `json:"catalogID" bson:"catalogID"`
	Status      int32    `json:"status" bson:"status"`
	Thumbnails  []string `json:"thumbnails" bson:"thumbnails"`
}

// New a new product instance
func New() Product {
	p := Product{}
	return p
}

// CreateProductRequest struct
type CreateProductRequest struct {
	Product Product `json:"product"`
}

// CreateProductResponse ...
type CreateProductResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// UploadProductRequest struct
type UploadProductRequest struct {
	Body []byte
	Md5  string
	Name string
}

// UploadProductResponse ...
type UploadProductResponse struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// Failer is an interface that should be implemented by response types.
// Response encoders can check if responses are Failer, and if so if they've
// failed, and if so encode them using a separate write path based on the error.
type Failer interface {
	Failed() error
}

// GetProductsRequest collects the request parameters for the GetProducts method.
type GetProductsRequest struct {
	A, B int64
}

// GetProductsResponse collects the response values for the GetProducts method.
type GetProductsResponse struct {
	V   int64 `json:"v"`
	Err error `json:"-"` // should be intercepted by Failed/errorEncoder
}

// Failed implements Failer.
func (r GetProductsResponse) Failed() error { return r.Err }
