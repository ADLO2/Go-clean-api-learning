package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	//"github.com/google/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
)

// News UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// News UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository, logger logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo, logger: logger}
}

// Create news
func (u *authUC) Register(ctx context.Context, registerRequest *models.RegisterRequest) (*models.User, error){
	newUser := &models.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Name: registerRequest.Name,
	}
	
	hash := md5.Sum([]byte(registerRequest.Password))
	hashedPassword := hex.EncodeToString(hash[:])
	registerRequest.Password = hashedPassword

	status, err := u.authRepo.CreateNewUser(ctx, registerRequest)
	if err != nil{
		return nil, err
	}

	if status != "register successful" {
		return nil, errors.BadQueryParams
	}

	return newUser, nil
}

// Get news by id
func (u *authUC) Login(ctx context.Context, loginRequest *models.LoginRequest) (string, string, error) {
	hash := md5.Sum([]byte(loginRequest.Password))
	hashedPassword := hex.EncodeToString(hash[:])
	loginRequest.Password = hashedPassword
	accessToken, refreshToken, key, err := u.authRepo.LoginAsUser(ctx, loginRequest)
	if err != nil {
        return "Login fail", "Login fail", err
    }

	tokens := &models.Token{
		AccessToken: accessToken,
	}

	err = u.redisRepo.SetJWTToken(ctx, key, 60, tokens)
	if err != nil {
        return "Login fail", "Login fail", err
    }
	
	return accessToken, refreshToken, nil
}

// Delete news
func (u *authUC) Logout(ctx context.Context, key string) error {
	return nil
}

func (u *authUC) GetNewAccessToken(ctx context.Context, key string, refreshToken string) (string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(key)
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	
	now := time.Now()
	accessToken := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expiredAt, ok := claims["ExpiredAt"].(string)
		if !ok {
			return "", errors.InvalidJWTClaims 
		}
		expiredTime, err := time.Parse(time.RFC3339Nano, expiredAt)
		if err != nil {
			return "", err
		}
		refreshKey, ok := claims["key"].(string)
		if !ok {
			return "", errors.InvalidJWTClaims 
		}
		if (now.Before(expiredTime)) {
			newLoginRequest := &models.LoginRequest{
				Username: claims["Username"].(string),
				Password: claims["Password"].(string),
			}
			_accessToken, er := u.authRepo.LoginWithRefreshToken(refreshKey, newLoginRequest)
			if er != nil {
				return "", er
			}
			fmt.Println(_accessToken)
			tokens := &models.Token{
				AccessToken: _accessToken,
			}
			er = u.redisRepo.SetJWTToken(ctx, refreshKey, 30, tokens)
			if er != nil {
				return "", er
			}
			accessToken = _accessToken
		} else {
			return "login expired, please login again", errors.NewErrorWithMessage(401, "Please login again", "login session expired")
		}
	}
	
	return accessToken, nil
}

