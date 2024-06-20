//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package auth

import (
	"context"

	// "github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	// "github.com/thienkb1123/go-clean-arch/pkg/utils"
)

// News use case
type UseCase interface {
	Register(ctx context.Context, registerRequest *models.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, loginRequest *models.LoginRequest) (string, string, error)
	GetNewAccessToken(ctx context.Context, key string, refreshToken string) (string, error)
	Logout(ctx context.Context, key string) error
}
