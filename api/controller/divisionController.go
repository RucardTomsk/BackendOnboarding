package controller

import (
	"context"
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/api/model"
	"github.com/RucardTomsk/BackendOnboarding/internal/api"
	"github.com/RucardTomsk/BackendOnboarding/internal/api/middleware"
	"github.com/RucardTomsk/BackendOnboarding/internal/domain/base"
	"github.com/RucardTomsk/BackendOnboarding/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type DivisionController struct {
	logger  *zap.Logger
	service *service.DivisionService
}

func NewDivisionController(
	logger *zap.Logger,
	service *service.DivisionService) *DivisionController {
	return &DivisionController{
		logger:  logger,
		service: service,
	}
}

// CreateDivision
// @Summary      Create Division
// @Description  Create Division
// @Tags         division
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateDivisionRequest true "User request"
// @Success      200  {object}  base.ResponseOKWithGUID "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /division/create [post]
func (a *DivisionController) CreateDivision(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.CreateDivisionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	id, serviceErr := a.service.Create(&payload, context.TODO())
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	log.Info(fmt.Sprintf("entity saved to db with id: %s", id.String()))

	c.JSON(http.StatusOK, base.ResponseOKWithGUID{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
		GUID:       *id,
	})
}

// GetDivisions
// @Summary      Get Divisions
// @Description  Get Divisions
// @Tags         division
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GetDivisions "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /division/all [get]
func (a *DivisionController) GetDivisions(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	divisions, serviceErr := a.service.Get(context.TODO())
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, model.GetDivisions{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c)},
		Divisions: divisions,
	})
}
