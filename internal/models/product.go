package models

import (
	"time"

	"github.com/google/uuid"
)

// Product base model
type Product struct {
	ProductID    int `json:"product_id" gorm:"column:product_id" validate:"omitempty"`
	AuthorID  uuid.UUID `json:"author_id,omitempty" gorm:"column:author_id" validate:"required"`
	Name     string    `json:"name" gorm:"column:name" validate:"required"`
	Description   string    `json:"description" gorm:"column:description" validate:"required"`
	Price int `json:"price" gorm:"column:price" validate:"required"`
	ImageURL  *string   `json:"image_url,omitempty" gorm:"column:image_url"`
	Category  *string   `json:"category,omitempty" gorm:"column:category"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

func (*Product) TableName() string {
	return "product"
}

// All Product response
type ProductList struct {
	TotalCount int64   `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Product       []*Product `json:"product"`
}

// Product base
type ProductBase struct {
	ProductID    int `json:"product_id" gorm:"column:product_id" validate:"omitempty"`
	AuthorID  uuid.UUID `json:"author_id,omitempty" gorm:"column:author_id" validate:"required"`
	Name     string    `json:"name" gorm:"column:name" validate:"required"`
	Description   string    `json:"description" gorm:"column:description" validate:"required"`
	Price int `json:"price" gorm:"column:price" validate:"required"`
	ImageURL  *string   `json:"image_url,omitempty" gorm:"column:image_url"`
	Category  *string   `json:"category,omitempty" gorm:"column:category"`
}

type ProductComment struct {
	CommentID uuid.UUID
	Product ProductBase
	User UserBase
	Content ContentBase
}

type CommentInput struct {
	AccessToken string 	`json:"accessToken,omitempty"`
	ProductID int		`json:"productID,omitempty"`
	Text string			`json:"text,omitempty"`
	Rating int			`json:"rating,omitempty"`
	ImageURL  []string 	`json:"imageURL,omitempty"`
}