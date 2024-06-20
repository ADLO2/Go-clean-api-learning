package http

import (
	"github.com/gin-gonic/gin"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/internal/security"
)

// Map news routes
func MapAuthRoutes(authGroup *gin.RouterGroup, h auth.Handlers, s *security.SecurityManager) {
	authGroup.Use(s.RateLimitCheck("authModule"))
	authGroup.POST("/register", h.Register)
	authGroup.POST("/login", h.Login)
	authGroup.GET("/refreshAToken", h.RefreshAccessToken)
}
