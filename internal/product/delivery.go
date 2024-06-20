package product

import (
	"github.com/gin-gonic/gin"
)

// News HTTP Handlers interface
type Handlers interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	GetByID(c *gin.Context)
	Delete(c *gin.Context)
	GetProduct(c *gin.Context)
	Comment(c *gin.Context)
	GetCommentByID(c *gin.Context)
	GetComments(c *gin.Context)
	DeleteComment(c *gin.Context)
}
