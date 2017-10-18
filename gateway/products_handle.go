package gateway

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	product_client "github.com/laidingqing/dabanshan/client/products"
	"github.com/laidingqing/dabanshan/proto/product"
)

func RegisterProducts(router *gin.RouterGroup) {
	r := router.Group("/products")
	r.GET("/get_products", GetProducts)
}

func GetProducts(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	size, err := strconv.ParseInt(c.Query("size"), 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	req := &product.GetProductsRequest{userID, size}
	resp, err := product_client.GetClient().GetProducts(context.Background(), req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.IndentedJSON(http.StatusOK, resp)
}
