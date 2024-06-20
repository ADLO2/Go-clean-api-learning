package repository

import (
	"context"
	"time"

	//"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/cache/redis"
	"github.com/thienkb1123/go-clean-arch/pkg/converter"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
)

// News redis repository
type authRedisRepo struct {
	rdb redis.Client
}

// News redis repository constructor
func NewAuthRedisRepo(rdb redis.Client) auth.RedisRepository {
	return &authRedisRepo{rdb: rdb}
}

// Get new by id
func (n *authRedisRepo) GetJWTToken(ctx context.Context, key string) (string, error) {
	authBytes, err := n.rdb.Get(ctx, key)
	if err != nil {
		return "nil", errors.WithMessage(err, "authRedisRepo.GetAuthByIDCtx.redisClusterClient.Get")
	}
	authBase := &models.Token{}
	if err = converter.BytesToAny(authBytes, authBase); err != nil {
		return "nil", errors.WithMessage(err, "authRedisRepo.GetAuthByIDCtx.json.Unmarshal")
	}

	return authBase.AccessToken, nil
}

// Cache news item
func (n *authRedisRepo) SetJWTToken(ctx context.Context, key string, seconds int, tokens *models.Token) error {
	authBytes, err := converter.AnyToBytes(tokens)
	if err != nil {
		return errors.WithMessage(err, "authRedisRepo.SetAuthCtx.json.Marshal")
	}
	if err = n.rdb.Set(ctx, key, authBytes, time.Second*time.Duration(seconds)); err != nil {
		return errors.WithMessage(err, "authRedisRepo.SetAuthCtx.redisClusterClient.Set")
	}
	return nil
}

// Delete new item from cache
func (n *authRedisRepo) DeleteJWTToken(ctx context.Context, key string) error {
	if err := n.rdb.Del(ctx, key); err != nil {
		return errors.WithMessage(err, "authRedisRepo.DeleteAuthCtx.redisClusterClient.Del")
	}
	return nil
}
