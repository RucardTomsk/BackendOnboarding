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
	"github.com/google/uuid"
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

// RegistrationUser
// @Summary      Registration User
// @Description  Registration User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateUserRequest true "User request"
// @Success      200  {object}  base.ResponseOKWithJWT "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/register [post]
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

// LoginUser
// @Summary      Login User
// @Description  Login User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        payload body model.LoginUserRequest true "User request"
// @Success      200  {object}  base.ResponseOKWithJWT "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/login [post]
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

// GetUserInfo
// @Summary      Get User Info
// @Description  Get User Info
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.UserInfoResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/info [get]
func (a *UserController) GetUserInfo(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	userIDAny, _ := c.Get(middleware.UserIDKey)
	userID := userIDAny.(string)
	userGuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	email, serviceErr := a.service.GetEmail(userGuid, context.TODO())
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, model.UserInfoResponse{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c),
		},
		Email: email,
	})
}

// GetUsers
// @Summary      Get Users
// @Description  Get Users
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GetUsersRequest "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/all [get]
func (a *UserController) GetUsers(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	users, serviceErr := a.service.Get(context.TODO())
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, model.GetUsersRequest{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c),
		},
		Users: users,
	})
}
