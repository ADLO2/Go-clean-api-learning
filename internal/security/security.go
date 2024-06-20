package security

import (
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/pkg/cache/redis"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
)

type SecurityManager struct {
	cfg     *config.Config
	origins []string
	logger  logger.Logger
	SCache  redis.Client
}

// Middleware manager constructor
func NewSecurityManager(cfg *config.Config, origins []string, SCache redis.Client, logger logger.Logger) *SecurityManager {
	return &SecurityManager{cfg: cfg, origins: origins, SCache: SCache, logger: logger}
}
