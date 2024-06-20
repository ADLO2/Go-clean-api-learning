package repository

import (
	"context"
	"time"

	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/product"
	"github.com/thienkb1123/go-clean-arch/pkg/cache/redis"
	"github.com/thienkb1123/go-clean-arch/pkg/converter"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
)

// News redis repository
type productRedisRepo struct {
	rdb redis.Client
}

// News redis repository constructor
func NewProductRedisRepo(rdb redis.Client) product.RedisRepository {
	return &productRedisRepo{rdb: rdb}
}

// Get new by id
func (n *productRedisRepo) GetProductByIDCtx(ctx context.Context, key string) (*models.ProductBase, error) {
	productBytes, err := n.rdb.Get(ctx, key)
	if err != nil {
		return nil, errors.WithMessage(err, "productRedisRepo.GetProductByIDCtx.redisClusterClient.Get")
	}
	productBase := &models.ProductBase{}
	if err = converter.BytesToAny(productBytes, productBase); err != nil {
		return nil, errors.WithMessage(err, "productRedisRepo.GetProductByIDCtx.json.Unmarshal")
	}

	return productBase, nil
}

// Cache news item
func (n *productRedisRepo) SetProductCtx(ctx context.Context, key string, seconds int, product *models.ProductBase) error {
	productBytes, err := converter.AnyToBytes(product)
	if err != nil {
		return errors.WithMessage(err, "productRedisRepo.SetProductCtx.json.Marshal")
	}
	if err = n.rdb.Set(ctx, key, productBytes, time.Second*time.Duration(seconds)); err != nil {
		return errors.WithMessage(err, "productRedisRepo.SetProductCtx.redisClusterClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *productRedisRepo) DeleteProductCtx(ctx context.Context, key string) error {
	if err := n.rdb.Del(ctx, key); err != nil {
		return errors.WithMessage(err, "productRedisRepo.DeleteProductCtx.redisClusterClient.Del")
	}
	return nil
}
