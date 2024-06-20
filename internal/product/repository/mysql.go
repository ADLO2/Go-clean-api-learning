package repository

import (
	"context"

	//"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/product"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
	"gorm.io/gorm"
)

// News Repository
type productRepo struct {
	db *gorm.DB
}

// News repository constructor
func NewProductRepository(db *gorm.DB) product.Repository {
	return &productRepo{db: db}
}

// Create news
func (r *productRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := r.db.Model(&models.Product{}).Create(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Update news item
func (r *productRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	db := r.db.Model(&models.Product{})
	err := db.First(&product).Error
	if err != nil {
		return nil, err
	}

	err = db.Save(product).Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Get single news by id
func (r *productRepo) GetProductByID(ctx context.Context, productID int) (*models.ProductBase, error) {
	result := &models.ProductBase{}
	err := r.db.Model(&models.Product{}).Where("product_id", productID).First(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete news by id
func (r *productRepo) Delete(ctx context.Context, productID int) error {
	err := r.db.Model(&models.Product{}).
		Where("product_id", productID).
		Delete(&models.Product{}).Error
	return err
}

// Get news
func (r *productRepo) GetProduct(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error) {
	totalCount := int64(0)
	db := r.db.WithContext(ctx).Model(&models.Product{})
	db.Count(&totalCount)

	if totalCount == 0 {
		return &models.ProductList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Product:    make([]*models.Product, 0),
		}, nil
	}

	productList := make([]*models.Product, 0, pq.GetSize())
	err := db.Offset(pq.GetOffset()).Limit(pq.GetLimit()).Find(&productList).Error
	if err != nil {
		return nil, err
	}
	return &models.ProductList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Product:    productList,
	}, nil
}
