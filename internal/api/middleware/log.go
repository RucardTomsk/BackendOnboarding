package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
)

// SetRequestLogging adds request logging to middleware chain.
// This function simply logs HTTP URL and body.
func SetRequestLogging(logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.With(zap.String("trackingID", GetTrackingId(c))).
			Info(fmt.Sprintf("%s %s %s (%s, %s)",
				c.Request.Proto, c.Request.Method, c.Request.RequestURI,
				c.ClientIP(), c.Request.UserAgent()))

		data, _ := io.ReadAll(c.Request.Body)
		logger.Debug(string(data))

		defer func() {
			logger.With(zap.String("trackingID", GetTrackingId(c))).
				Info(fmt.Sprintf("respond with %d", c.Writer.Status()))
		}()

		// and now set a new body, which will simulate the same data we read
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		c.Next()
	}
}
