//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package product

import (
	"context"

	//"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
)

// News use case
type UseCase interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProductByID(ctx context.Context, productID int) (*models.ProductBase, error)
	Delete(ctx context.Context, productID int) error
	GetProduct(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error)
	Comment(ctx context.Context, user *models.UserBase, productID int, content *models.ContentBase) (*models.ProductComment, error)
}
