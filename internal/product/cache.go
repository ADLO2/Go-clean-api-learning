//go:generate mockgen -source redis_repository.go -destination mock/redis_repository_mock.go -package mock
package product

import (
	"context"

	"github.com/thienkb1123/go-clean-arch/internal/models"
)

// News redis repository
type RedisRepository interface {
	GetProductByIDCtx(ctx context.Context, key string) (*models.ProductBase, error)
	SetProductCtx(ctx context.Context, key string, seconds int, product *models.ProductBase) error
	DeleteProductCtx(ctx context.Context, key string) error
}
