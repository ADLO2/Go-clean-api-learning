package middleware

import (
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
	RCache  auth.RedisRepository
}

// Middleware manager constructor
func NewMiddlewareManager(cfg *config.Config, origins []string, RCache auth.RedisRepository, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, origins: origins, RCache: RCache, logger: logger}
}
