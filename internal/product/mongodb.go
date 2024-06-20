package product

import (
	"context"

	"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
)

// News Repository
type MongoRepository interface {
	Create(ctx context.Context, comment *models.ProductComment) error
	GetProductCommentByID(ctx context.Context, commentID uuid.UUID) (*models.ProductComment, error)
	DeleteProductComment(ctx context.Context, commentID uuid.UUID) error
	GetProductComments(ctx context.Context, productID int) (*models.ProductComment, error)
}
