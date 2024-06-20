package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID `json:"user_id" validate:"omitempty"`
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty"`
	Name string `json:"name" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty"`
	Name string `json:"name" validate:"omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty"`
}

type JWTPayload struct {
	UserID uuid.UUID `json:"user_id" validate:"omitempty"`
	Name string `json:"name" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	ExpiredAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}

type Token struct {
	AccessToken string 	`json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserBase struct {
	UserID uuid.UUID `json:"user_id" validate:"omitempty"`
	Username string 	`json:"Username"`
}