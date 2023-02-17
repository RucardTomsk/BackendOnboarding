package api

import (
	"github.com/RucardTomsk/BackendOnboarding/internal/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// EnrichLogger returns logger with retrieved fields from gin.Context:
// trackingID, sessionID and operackend_onboardingation name.
func EnrichLogger(logger *zap.Logger, c *gin.Context) *zap.Logger {
	return logger.With(
		zap.String("trackingID", middleware.GetTrackingId(c)),
		zap.String("operation", middleware.GetOperationName(c)),
		zap.String("sessionID", middleware.GetSessionId(c)),
	)
}
