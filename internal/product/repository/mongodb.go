package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/product"

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// News Repository
type MongoRepository struct {
	cli *mongo.Client
}

// News repository constructor
func NewProductCommentRepo(cli *mongo.Client) product.MongoRepository {
	return &MongoRepository{cli: cli}
}

func (m *MongoRepository) Create(ctx context.Context, comment *models.ProductComment) error{ 
	collection := m.cli.Database("comments").Collection("product_comments")
	_, err := collection.InsertOne(context.Background(), comment)
	return err
}

func (m *MongoRepository) GetProductCommentByID(ctx context.Context, commentID uuid.UUID) (*models.ProductComment, error){
	return nil, nil
}

func (m *MongoRepository) DeleteProductComment(ctx context.Context, commentID uuid.UUID) error {
	return nil
}

func (m *MongoRepository) GetProductComments(ctx context.Context, productID int) (*models.ProductComment, error){
	return nil, nil
}
