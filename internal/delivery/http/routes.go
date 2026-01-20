package http

import (
	"fmt"
	"net/http"
	"time"
	"worker-service/config"
	"worker-service/internal/controller"
	"worker-service/internal/dto"
	"worker-service/internal/middleware"
	"worker-service/internal/pkg/redis"
	"worker-service/internal/repository"
	"worker-service/internal/repository/unitofwork"
	"worker-service/internal/services"
	"worker-service/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	ApiKeyMiddleware    gin.HandlerFunc
	RateLimitMiddleware gin.HandlerFunc
	EmailController     controller.EmailController
}

func InitRoutes(db *gorm.DB) *gin.Engine {
	return setupRoutes(*initHandler(db))
}

func initHandler(db *gorm.DB) *Handlers {
	// Validator
	appConfig := config.New()
	uow := unitofwork.NewUoW(db)

	// Auth
	apiKeyRepository := repository.NewApiKeyRepository(db)
	authUsecase := usecase.NewAuthUsecase(apiKeyRepository)
	apiKeyMiddleware := middleware.APIKeyMiddleware(authUsecase)

	// Rate Limit
	cache := redis.NewRedisClient[services.TokenBucket](*appConfig, "rate_limit:mail_service", 1*time.Minute)
	rateLimitService := services.NewRateLimiter(cache)
	rateLimitMiddleware := middleware.RateLimitterMiddleware(*rateLimitService)

	// Email
	emailHistoryRepository := repository.NewEmailHistoryRepository(db)
	emailService := services.NewEmailService(*appConfig)
	redisClient := redis.NewRedisClient[dto.EmailTask](*appConfig, "email_queue", 0)
	emailUsecase := usecase.NewEmailUsecase(appConfig, emailHistoryRepository, uow, emailService, redisClient)
	emailController := controller.NewEmailController(emailUsecase)

	return &Handlers{
		EmailController:     emailController,
		ApiKeyMiddleware:    apiKeyMiddleware,
		RateLimitMiddleware: rateLimitMiddleware,
	}
}

func setupRoutes(handler Handlers) *gin.Engine {
	route := gin.Default()

	// CORS handler
	route.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{"POST", "PUT", "PATCH", "DELETE", "GET", "OPTIONS", "TRACE", "CONNECT"},
		AllowHeaders:     []string{"Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Origin", "Content-Type", "Content-Length", "Date", "origin", "Origins", "x-requested-with", "access-control-allow-methods", "access-control-allow-credentials", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
	}))

	// Use panic recover middleware
	route.NoRoute(noRoute)

	// API group
	api := route.Group("/api/v1")
	api.GET("/", index)

	// Email
	api.Use(handler.ApiKeyMiddleware)
	api.GET(controller.EmailPath, handler.EmailController.ListEmail)
	api.GET(controller.EmailByIdPath, handler.EmailController.ListEmailByID)
	api.POST(controller.EmailSendBulkPath, handler.RateLimitMiddleware, handler.EmailController.SendEmail)
	api.POST(controller.EmailRetryPath, handler.RateLimitMiddleware, handler.EmailController.RetryEmail)

	return route
}

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "welcome to kenneth's email service",
	})
}

func noRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
	})
}
