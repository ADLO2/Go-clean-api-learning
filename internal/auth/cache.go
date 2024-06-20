//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package auth

import (
	"context"

	"github.com/thienkb1123/go-clean-arch/internal/models"
)

// News redis repository
type RedisRepository interface {
	GetJWTToken(ctx context.Context, key string) (string, error)
	SetJWTToken(ctx context.Context, key string, seconds int, tokens *models.Token) error
	DeleteJWTToken(ctx context.Context, key string) error
}
