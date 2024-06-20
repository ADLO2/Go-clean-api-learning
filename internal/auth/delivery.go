package auth

import (
	"github.com/gin-gonic/gin"
)

// News HTTP Handlers interface
type Handlers interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshAccessToken(c *gin.Context)
}
