package middleware

import (
	"github.com/RucardTomsk/BackendOnboarding/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	UserIDKey           = "UserID"
)

// AuthorizationCheck adds authorization check to middleware chain.
func AuthorizationCheck(service *service.UserService, logger zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.GetHeader(authorizationHeader)
		if header == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string
			}{Error: "ERROR MIDDLE AUTH"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string
			}{Error: "ERROR MIDDLE AUTH"})
			return
		}

		userGuid, err := service.ParseToken(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, struct {
				Error string
			}{Error: "ERROR MIDDLE AUTH"})
			return
		}

		c.Set(UserIDKey, userGuid)
		c.Next()
	}
}
