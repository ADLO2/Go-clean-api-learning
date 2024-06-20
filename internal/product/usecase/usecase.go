package usecase

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/product"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
)

const (
	basePrefix    = "api-product:"
	cacheDuration = 3600
)

// News UseCase
type productUC struct {
	cfg       *config.Config
	productRepo  product.Repository
	redisRepo product.RedisRepository
	mongoRepo product.MongoRepository
	logger    logger.Logger
}

// News UseCase constructor
func NewProductUseCase(cfg *config.Config, productRepo product.Repository, redisRepo product.RedisRepository, mongoRepo product.MongoRepository,logger logger.Logger) product.UseCase {
	return &productUC{cfg: cfg, productRepo: productRepo, redisRepo: redisRepo, mongoRepo: mongoRepo, logger: logger}
}

// Create news
func (u *productUC) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	user, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		return nil, errors.NewUnauthorizedError(errors.WithMessage(err, "newsUC.Create.GetUserFromCtx"))
	}

	product.AuthorID = user.UserID
	if err = utils.ValidateStruct(ctx, product); err != nil {
		return nil, errors.NewBadRequestError(errors.WithMessage(err, "newsUC.Create.ValidateStruct"))
	}

	n, err := u.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return n, err
}

// Update news item
func (u *productUC) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	productByID, err := u.productRepo.GetProductByID(ctx, product.ProductID)
	if err != nil {
		return nil, err
	}



	if err = utils.ValidateIsOwner(ctx, productByID.AuthorID.String(), u.logger); err != nil {
		return nil, errors.NewError(http.StatusForbidden, errors.ErrForbidden, errors.WithMessage(err, "productUC.Update.ValidateIsOwner"))
	}

	updatedUser, err := u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.DeleteProductCtx(ctx, u.getKeyWithPrefix(strconv.Itoa(product.ProductID))); err != nil {
		u.logger.Errorf(ctx, "productUC.Update.DeleteProductCtx: %v", err)
	}

	return updatedUser, nil
}

// Get news by id
func (u *productUC) GetProductByID(ctx context.Context, productID int) (*models.ProductBase, error) {
	productBase, err := u.redisRepo.GetProductByIDCtx(ctx, u.getKeyWithPrefix(strconv.Itoa(productID)))
	if err != nil {
		u.logger.Errorf(ctx, "productUC.GetNewsByID.GetNewsByIDCtx: %v", err)
	}
	if productBase != nil {
		return productBase, nil
	}

	n, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetProductCtx(ctx, u.getKeyWithPrefix(strconv.Itoa(productID)), cacheDuration, n); err != nil {
		u.logger.Errorf(ctx, "productUC.GetProductByID.SetNewsCtx: %s", err)
	}

	return n, nil
}

// Delete news
func (u *productUC) Delete(ctx context.Context, productID int) error {
	productByID, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	if err = utils.ValidateIsOwner(ctx, productByID.AuthorID.String(), u.logger); err != nil {
		return errors.NewError(http.StatusForbidden, errors.ErrForbidden, errors.WithMessage(err, "productUC.Delete.ValidateIsOwner"))
	}

	if err = u.productRepo.Delete(ctx, productID); err != nil {
		return err
	}

	if err = u.redisRepo.DeleteProductCtx(ctx, u.getKeyWithPrefix(strconv.Itoa(productID))); err != nil {
		u.logger.Errorf(ctx, "productUC.Delete.DeleteProductCtx: %v", err)
	}

	return nil
}

// Get news
func (u *productUC) GetProduct(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error) {
	results, err := u.productRepo.GetProduct(ctx, pq)
	if err != nil {
		u.logger.Error(ctx, err)
	}
	return results, err
}

func (u *productUC) getKeyWithPrefix(productID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, productID)
}

func (u *productUC) Comment(ctx context.Context, user *models.UserBase, productID int, content *models.ContentBase) (*models.ProductComment, error){
	product, err := u.productRepo.GetProductByID(ctx, productID)
	if err != nil {
		u.logger.Error(ctx, err)
	}
	comment := &models.ProductComment{
		CommentID: uuid.New(),
		User: *user,
		Product: *product,
		Content: *content,
	}
	err = u.mongoRepo.Create(ctx, comment)
	if err != nil {
		u.logger.Error(ctx, err)
	}

	return comment, err
}