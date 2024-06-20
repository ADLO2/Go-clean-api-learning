package http

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/google/uuid"
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/internal/product"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
	"github.com/thienkb1123/go-clean-arch/pkg/response"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
)

type productHandlers struct {
	cfg    *config.Config
	productUC product.UseCase
	logger logger.Logger
}

func NewProductHandlers(cfg *config.Config, productUC product.UseCase, logger logger.Logger) product.Handlers {
	return &productHandlers{cfg: cfg, productUC: productUC, logger: logger}
}

func (h productHandlers) Create(c *gin.Context) {

	n := &models.Product{}
	if err := c.Bind(n); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	createdProduct, err := h.productUC.Create(ctx, n)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, createdProduct)
}

func (h productHandlers) Update(c *gin.Context) {
	productUUID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	n := &models.Product{}
	if err := c.Bind(n); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}
	n.ProductID = productUUID

	ctx := c.Request.Context()
	updatedProduct, err := h.productUC.Update(ctx, n)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, updatedProduct)
}

func (h productHandlers) GetByID(c *gin.Context) {
	productUUID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	productByID, err := h.productUC.GetProductByID(ctx, productUUID)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, productByID)
}

func (h productHandlers) Delete(c *gin.Context) {
	productUUID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	if err = h.productUC.Delete(ctx, productUUID); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithNoContent(c)
}

func (h productHandlers) GetProduct(c *gin.Context) {
	pq, err := utils.GetPaginationFromCtx(c)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	ctx := c.Request.Context()
	productList, err := h.productUC.GetProduct(ctx, pq)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, productList)
}

func (h productHandlers) Comment(c *gin.Context) {
	n := &models.CommentInput{}

	if err := c.ShouldBind(n); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	token, err := jwt.Parse(n.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(h.cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	if !token.Valid {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, errors.InvalidJWTToken)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["UserID"].(string)
		if !ok {
			utils.LogResponseError(c, h.logger, err)
			response.WithError(c, errors.InvalidJWTToken)
			return
		}
		username, ok := claims["Username"].(string)
		fmt.Println(username)
		if !ok {
			utils.LogResponseError(c, h.logger, err)
			response.WithError(c, errors.InvalidJWTToken)
			return
		}

		userUUID, _ := uuid.Parse(userID)
		user :=&models.UserBase{
			UserID:  userUUID,
			Username:  username,
		}

		content := &models.ContentBase{
			Text: n.Text,
			Rating: n.Rating,
			ImageURL: n.ImageURL,
		}

		ctx := c.Request.Context()
		productComment, err := h.productUC.Comment(ctx, user, n.ProductID, content)
		if err != nil {
			utils.LogResponseError(c, h.logger, err)
			response.WithError(c, err)
			return
		}
		response.WithOK(c, productComment)
	} else {
		utils.LogResponseError(c, h.logger, errors.InvalidJWTToken)
		response.WithError(c, errors.InvalidJWTToken)
		return
	}
}

func (h productHandlers) GetCommentByID(c *gin.Context) {
	
}

func (h productHandlers) GetComments(c *gin.Context) {
	
}

func (h productHandlers) DeleteComment(c *gin.Context) {
	
}