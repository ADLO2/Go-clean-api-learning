package repository

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/internal/models"

	// "github.com/thienkb1123/go-clean-arch/pkg/utils"
	"gorm.io/gorm"
)

// News Repository
type productRepo struct {
	db *gorm.DB
}

func generateRandomString(length int) string{
	ran_str := make([]byte, length)
	for i := 0; i < length; i++ { 
        ran_str[i] = byte(65 + rand.Intn(25))
    } 
  
    // Displaying the random string 
    str := string(ran_str)
	fmt.Println(str)
	return str
}

// News repository constructor
func NewAuthRepository(db *gorm.DB) auth.Repository {
	return &productRepo{db: db}
}

// Create news
func (r *productRepo) CreateNewUser(ctx context.Context, registerRequest *models.RegisterRequest) (string, error) {
	newUser := &models.User{
		UserID: uuid.New(),
		Username: registerRequest.Username,
		Name: registerRequest.Name,
		Password: registerRequest.Password,
	}

	err := r.db.Model(&models.User{}).Create(newUser).Error
	if err != nil {
		return "register fail", err
	}

	return "register successful", nil
}

// Get single news by id
func (r *productRepo) LoginAsUser(ctx context.Context, loginRequest *models.LoginRequest) (string, string, string, error) {
	result := &models.User{}
	err := r.db.Model(&models.User{}).Where("username = ? and password = ?", loginRequest.Username, loginRequest.Password).First(&result).Error
	if err != nil {
		return "Login fail", "Login fail", "Login fail", err
	}
	err = godotenv.Load(".env")

    if err != nil {
        return "Login fail", "Login fail", "Login fail", err
    }


	rTokenKey := "authKey:" + generateRandomString(8)
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key": rTokenKey,
		"UserID": result.UserID,
		"Username": result.Username,
		"Password": result.Password,
		"CreatedAt": time.Now(),
		"ExpiredAt": time.Now().Add(time.Minute * 5),
	})
	randomize := generateRandomString(8)
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key": rTokenKey,
		"UserID": result.UserID,
		"Username": result.Username,
		"randomize": randomize,
	})

	accessToken, err := jwtAccessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
        return "Login fail", "Login fail", "Login fail", err
    }
	refreshToken, err := jwtRefreshToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
        return "Login fail", "Login fail", "Login fail", err
    }

	return accessToken, refreshToken, rTokenKey, nil
}

func (r *productRepo) LoginWithRefreshToken(accessTokenKey string, loginRequest *models.LoginRequest) (string, error) {
	result := &models.User{}
	err := r.db.Model(&models.User{}).Where("username = ? and password = ?", loginRequest.Username, loginRequest.Password).First(&result).Error
	if err != nil {
		return "Login fail", err
	}
	err = godotenv.Load(".env")

    if err != nil {
        return "Login fail", err
    }
	randomize := generateRandomString(8)
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"key": accessTokenKey,
		"UserID": result.UserID,
		"randomize": randomize,
	})

	accessToken, err := jwtAccessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
        return "Login fail", err
    }

	return accessToken, nil
}