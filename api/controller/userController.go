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

type UserController struct {
	logger  *zap.Logger
	service *service.UserService
}

func NewUserController(logger *zap.Logger, service *service.UserService) *UserController {
	return &UserController{
		logger:  logger,
		service: service,
	}
}

func (a *UserController) RegistrationUser(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.CreateUserRequest
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

	token, serviceErr := a.service.GenerateToken(
		&model.GenerateTokenRequest{
			Email:    payload.Email,
			Password: payload.Password,
		}, context.TODO())

	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithJWT{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
		GUID:       *id,
		JWT:        token.Value,
	})
}

func (a *UserController) LoginUser(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.LoginUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	token, serviceErr := a.service.GenerateToken(
		&model.GenerateTokenRequest{
			Email:    payload.Email,
			Password: payload.Password,
		}, context.TODO())

	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithJWT{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
		JWT:        token.Value,
	})
}
