package server

import (
	"context"

	"github.com/gin-contrib/requestid"
	authHttp "github.com/thienkb1123/go-clean-arch/internal/auth/delivery/http"
	authRepository "github.com/thienkb1123/go-clean-arch/internal/auth/repository"
	authUseCase "github.com/thienkb1123/go-clean-arch/internal/auth/usecase"
	apiMiddlewares "github.com/thienkb1123/go-clean-arch/internal/middleware"
	newsHttp "github.com/thienkb1123/go-clean-arch/internal/news/delivery/http"
	newsRepository "github.com/thienkb1123/go-clean-arch/internal/news/repository"
	newsUseCase "github.com/thienkb1123/go-clean-arch/internal/news/usecase"
	productHttp "github.com/thienkb1123/go-clean-arch/internal/product/delivery/http"
	productRepository "github.com/thienkb1123/go-clean-arch/internal/product/repository"
	productUseCase "github.com/thienkb1123/go-clean-arch/internal/product/usecase"
	security "github.com/thienkb1123/go-clean-arch/internal/security"
	"github.com/thienkb1123/go-clean-arch/pkg/metric"
)

// Map Server Handlers
func (s *Server) MapHandlers() error {
	ctx := context.Background()
	
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf(ctx, "CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		ctx,
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	metrics.SetSkipPath([]string{"readiness"})

	// Init repositories
	nRepo := newsRepository.NewNewsRepository(s.db)
	newsRedisRepo := newsRepository.NewNewsRedisRepo(s.redis)

	pRepo := productRepository.NewProductRepository(s.db)
	productRedisRepo := productRepository.NewProductRedisRepo(s.redis)
	pMongoRepo := productRepository.NewProductCommentRepo(s.mongodb)
	
	authRepo := authRepository.NewAuthRepository(s.db)
	authRedisRepo := authRepository.NewAuthRedisRepo(s.redis)

	// Init useCases
	newsUC := newsUseCase.NewNewsUseCase(s.cfg, nRepo, newsRedisRepo, s.logger)
	productUC := productUseCase.NewProductUseCase(s.cfg, pRepo, productRedisRepo, pMongoRepo, s.logger)
	authUC := authUseCase.NewAuthUseCase(s.cfg, authRepo, authRedisRepo, s.logger)

	// Init handlers
	newsHandlers := newsHttp.NewNewsHandlers(s.cfg, newsUC, s.logger)
	productHandlers := productHttp.NewProductHandlers(s.cfg, productUC, s.logger)
	authHandlers := authHttp.NewAuthHandlers(s.cfg, authUC, s.logger)


	mw := apiMiddlewares.NewMiddlewareManager(s.cfg, []string{"*"}, authRedisRepo,s.logger)
	secureManager := security.NewSecurityManager(s.cfg, []string{"*"}, s.redis, s.logger)


	s.gin.Use(requestid.New())
	s.gin.Use(mw.MetricsMiddleware(metrics))
	s.gin.Use(mw.LoggerMiddleware(s.logger))

	v1 := s.gin.Group("/api/v1")
	newsGroup := v1.Group("/news")

	v2 := s.gin.Group("/api/v2")
	productGroup := v2.Group("/product")

	v3 := s.gin.Group("/api/v3")
	authGroup := v3.Group("/auth")

	

	newsHttp.MapNewsRoutes(newsGroup, newsHandlers, mw)
	productHttp.MapProductRoutes(productGroup, productHandlers, mw)
	authHttp.MapAuthRoutes(authGroup, authHandlers, secureManager)

	


	return nil
}
