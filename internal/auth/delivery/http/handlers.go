package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/thienkb1123/go-clean-arch/config"
	"github.com/thienkb1123/go-clean-arch/internal/auth"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/errors"
	"github.com/thienkb1123/go-clean-arch/pkg/logger"
	"github.com/thienkb1123/go-clean-arch/pkg/response"
	"github.com/thienkb1123/go-clean-arch/pkg/utils"
)

// News handlers
type authHandlers struct {
	cfg    *config.Config
	userUC auth.UseCase
	logger logger.Logger
}

// NewNewsHandlers News handlers constructor
func NewAuthHandlers(cfg *config.Config, userUC auth.UseCase, logger logger.Logger) auth.Handlers {
	return &authHandlers{cfg: cfg, userUC: userUC, logger: logger}
}

func (h authHandlers) Register(c *gin.Context) {
	n := &models.RegisterRequest{}
	if err := c.ShouldBindBodyWith(n, binding.JSON); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	if n.Username == "" || n.Password == "" {
		utils.LogResponseError(c, h.logger, errors.Error{ErrMessage: "Empty username or password"})
		response.WithError(c, errors.Error{ErrMessage: "Empty username or password"})
		return
	}

	ctx := c.Request.Context()
	status, err := h.userUC.Register(ctx, n)
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	response.WithOK(c, status)
}


func (h authHandlers) Login(c *gin.Context) {
	n := &models.LoginRequest{}
	if err := c.ShouldBindBodyWith(n, binding.JSON); err != nil {
		fmt.Println("login")
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}

	if n.Username == "" || n.Password == "" {
		fmt.Println("login")
		utils.LogResponseError(c, h.logger, errors.Error{ErrMessage: "Empty username or password"})
		response.WithError(c, errors.Error{ErrMessage: "Empty username or password"})
		return
	}

	ctx := c.Request.Context()
	accessToken, refreshToken, err := h.userUC.Login(ctx, n)
	if err != nil {
		fmt.Println("login")
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}
	// create authKey: accessValue to redis


	result := map[string]string {"accessToken": accessToken, "refreshToken": refreshToken}
	response.WithOK(c, result)
}

func (h authHandlers) RefreshAccessToken(c *gin.Context){
	requestBody := map[string]string {
		"refreshToken": "",
	}
	if err := c.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}
	if requestBody["refreshToken"] == "" {
		response.WithError(c, errors.InvalidJWTToken)
		return
	}
	ctx := c.Request.Context()
	accessToken, err := h.userUC.GetNewAccessToken(ctx, h.cfg.Server.JwtSecretKey, requestBody["refreshToken"])
	if err != nil {
		utils.LogResponseError(c, h.logger, err)
		response.WithError(c, err)
		return
	}
	response.WithOK(c, accessToken)
}