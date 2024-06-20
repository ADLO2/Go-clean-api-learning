//go:generate mockgen -source mysql_repository.go -destination mock/mysql_repository_mock.go -package mock

package product

import (
	"context"

	// "github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
)

// News Repository
type Repository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductByID(ctx context.Context, productID int) (*models.ProductBase, error)
	Delete(ctx context.Context, productID int) error
	GetProduct(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error)
}
