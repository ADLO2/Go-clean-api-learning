package security

import (
	"context"
	"strconv"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis/v8"
	"github.com/thienkb1123/go-clean-arch/internal/models"
	"github.com/thienkb1123/go-clean-arch/pkg/converter"
	errorModule "github.com/thienkb1123/go-clean-arch/pkg/errors"
	"go.uber.org/zap"
)

func (s *SecurityManager) RateLimitCheck(moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		ctx := c.Request.Context()
		s.logger.Infof(ctx, "rate limit check security header %s", clientIP)
		err := s.checkLimit(ctx, moduleName + "IP", clientIP)
		if err != nil {
			s.logger.Error(ctx, "RateLimitCheck.CheckIP", zap.String("headerJWT", err.Error()))
			c.JSON(429, map[string]string{"code": "429", "message":"Too many request"})
			c.Abort()
			c.Next()
			return
		}

		n := &models.LoginRequest{}

		if err := c.ShouldBindBodyWith(n, binding.JSON); err != nil {
			s.logger.Error(ctx, "RateLimitCheck.CheckIP.GetUsernameFromCtx", zap.String("headerJWT", err.Error()))
			c.JSON(500, "Internel server error")
			c.Abort()
			c.Next()
			return
		}

		err = s.checkLimit(ctx, moduleName + "Username", n.Username)
		if err != nil {
			s.logger.Error(ctx, "RateLimitCheck.CheckUsername", zap.String("headerJWT", err.Error()))
			c.JSON(429, map[string]string{"code": "429", "message":"Too many request"})
			c.Abort()
			c.Next()
			return
		}
		c.Next()
	}
}

func (s *SecurityManager) checkLimit(ctx context.Context, moduleName string, k string) error {
	key := moduleName + "-" + k
	var limit int 
	if moduleName == "authModuleIP" {
		limit = 100
	} else if moduleName == "authModuleUsername" {
		limit = 5
	} else {
		return errorModule.NewErrorWithMessage(401, "UndefinedModule", nil)
	}
	limitRateBytes, err := s.SCache.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			_limit, e := converter.AnyToBytes(0)
			if e != nil {
				return errorModule.WithMessage(e, "securityRedis.CheckIP.SetIPLimit.json.Marshal")
			}
			
			e = s.SCache.Set(ctx, key, _limit, time.Hour * 24)
			if e != nil {
				return errorModule.WithMessage(e, "securityRedis.CheckIP.SetIPLimit")
			}
			return nil
		} else if err != nil {
			return errorModule.WithMessage(err, "securityRedis.CheckIP.GetIPlimit")
		}
	}

	rate, err := strconv.Atoi(string(limitRateBytes))
	if err != nil {
		return errorModule.WithMessage(err, "securityRedis.CheckIP.GetByKey.json.Unmarshal")
	}
	
	if rate >= limit {
		err = errors.New("securityRedis.CheckIP.LimitReach")
		e := errorModule.NewErrorWithMessage(429, "Too many request", err)
		return errorModule.WithMessage(e, "securityRedis.CheckIP.LimitReach")
	}
	s.SCache.Incr(ctx, key)
	return nil
}
