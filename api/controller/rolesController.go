package controller

import (
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/api/middleware"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/enum"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type RolesController struct {
	logger *zap.Logger
}

func NewRolesController(
	logger *zap.Logger) *RolesController {
	return &RolesController{
		logger: logger,
	}
}

// GetRoles
// @Summary      Get Roles
// @Description  Get Roles
// @Tags         roles
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GetRoles "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /roles/all [get]
func (a *RolesController) GetRoles(c *gin.Context) {

	roles := []string{
		enum.HR.String(),
		enum.ADMIN.String(),
		enum.BEGINNER.String(),
		enum.DIRECTOR.String(),
		enum.EMPLOYEE.String(),
	}

	c.JSON(http.StatusOK, model.GetRoles{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c),
		},
		Roles: roles,
	})
}
