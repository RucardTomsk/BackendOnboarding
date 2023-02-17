package middleware

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/common"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SetAccessControl adds access control to middleware chain.
// This function checks if user is able to perform certain action or not.
func SetAccessControl(cfg common.ServerConfig, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.With(zap.String("trackingID", GetTrackingId(c))).Debug("acl is not implemented")
		c.Next()
	}
}
