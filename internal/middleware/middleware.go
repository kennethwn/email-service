package middleware

import (
	"worker-service/internal/dto"
	"worker-service/internal/pkg/error_wrap"
	"worker-service/internal/services"
	"worker-service/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func APIKeyMiddleware(authUsecase usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request

		// Verify api key
		apiKey := req.Header.Get("X-API-KEY")
		if apiKey == "" {
			dto.WriteErrorResponseJSON(c, error_wrap.ErrApiKeyIsMissing)
			c.Abort()
			return
		}

		res, err := authUsecase.VerifyAPIKey(c, apiKey)
		if err != nil || !res.IsValid {
			dto.WriteErrorResponseJSON(c, error_wrap.ErrApiKeyIsInvalid)
			c.Abort()
			return
		}

		c.Set("api_key", res.APIKey)
		c.Set("rate_limit", res.ThresholdRateLimit)

		c.Next()
	}
}

func RateLimitterMiddleware(rateLimitter services.RateLimitter) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, isExist := c.Get("api_key")
		if !isExist {
			logrus.Error("error service is forbidden")
			dto.WriteErrorResponseJSON(c, error_wrap.ErrIPorServiceBlocked)
			c.Abort()
			return
		}

		token, isTokenExist := c.Get("rate_limit")
		if !isTokenExist {
			logrus.Error("error token not found")
			dto.WriteErrorResponseJSON(c, error_wrap.ErrIPorServiceBlocked)
			c.Abort()
			return
		}

		if err := rateLimitter.RefillToken(c, apiKey.(string), token.(int)); err != nil {
			logrus.Error("error refilling token: ", err)
			dto.WriteErrorResponseJSON(c, error_wrap.ErrIPorServiceBlocked)
			c.Abort()
			return
		}

		if !rateLimitter.Allow(c, apiKey.(string)) {
			logrus.Error("error rate limit")
			dto.WriteErrorResponseJSON(c, error_wrap.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}
