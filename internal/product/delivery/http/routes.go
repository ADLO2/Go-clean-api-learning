package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thienkb1123/go-clean-arch/internal/middleware"
	"github.com/thienkb1123/go-clean-arch/internal/product"
)

// Map news routes
func MapProductRoutes(productGroup *gin.RouterGroup, h product.Handlers, mw *middleware.MiddlewareManager) {
	productGroup.Use(mw.AuthJWTMiddleware())
	productGroup.POST("/create", h.Create)
	productGroup.PUT("/:productId", h.Update)
	productGroup.DELETE("/:productId", h.Delete)
	productGroup.GET("/:productId", h.GetByID)
	productGroup.GET("", h.GetProduct)
	productGroup.POST("/comment", h.Comment)
}
