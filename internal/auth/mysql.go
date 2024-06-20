//go:generate mockgen -source mysql_repository.go -destination mock/mysql_repository_mock.go -package mock

package auth

import (
	"context"

	// "github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	// "github.com/thienkb1123/go-clean-arch/pkg/utils"
)

// News Repository
type Repository interface {
	CreateNewUser(ctx context.Context, registerRequest *models.RegisterRequest) (string, error)
	LoginAsUser(ctx context.Context, loginRequest *models.LoginRequest) (string, string, string, error)
	LoginWithRefreshToken(accessTokenKey string, loginRequest *models.LoginRequest) (string, error)
}
