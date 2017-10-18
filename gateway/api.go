package gateway

import "github.com/gin-gonic/gin"

func Register(router *gin.Engine) {
	r := router.Group("/api")
	RegisterProducts(r)
}
