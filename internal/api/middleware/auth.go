package middleware

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetAuthorizationCheck adds authorization check to middleware chain.
func SetAuthorizationCheck(cfg common.ServerConfig, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.With(zap.String("trackingID", GetTrackingId(c))).Debug("auth is not implemented")
		c.Next()
	}
}
