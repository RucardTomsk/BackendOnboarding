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
	logger       *zap.Logger
	service      *service.UserService
	questService *service.QuestService
}

func NewUserController(logger *zap.Logger, service *service.UserService, questService *service.QuestService) *UserController {
	return &UserController{
		logger:       logger,
		service:      service,
		questService: questService,
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

func (a *UserController) Test(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)
	userIDAny, ok := c.Get(middleware.UserIDKey)
	if !ok {
		log.Warn("not JWT")
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}
	userID := userIDAny.(string)
	userGuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOKWithGUID{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
		GUID:       userGuid,
	})
}

// UpdateAbout
// @Summary      Update About
// @Description  Update About
// @Tags         user
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Authentication header"
// @Param        payload body model.UpdateAbout true "User request"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/about/update [post]
func (a *UserController) UpdateAbout(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	userIDAny, ok := c.Get(middleware.UserIDKey)
	if !ok {
		log.Warn("not JWT")
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}
	userID := userIDAny.(string)
	userGuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	var payload model.UpdateAbout
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	if serviceErr := a.service.UpdateAbout(&payload, userGuid, context.TODO()); err != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
	})
}

// GetUserQuest
// @Summary      Get User Quest
// @Description  Get User Quest
// @Tags         user
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Authentication header"
// @Success      200  {object}  model.AllQuestUserResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/quest [get]
func (a *UserController) GetUserQuest(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	userIDAny, ok := c.Get(middleware.UserIDKey)
	if !ok {
		log.Warn("not JWT")
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}
	userID := userIDAny.(string)
	userGuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	divQuest, serviceErr := a.questService.GetAllUserQuest(userGuid)
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, model.AllQuestUserResponse{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c),
		},
		DivQuests: divQuest,
	})
}

// GetUserInfo
// @Summary      Get User Info
// @Description  Get User Info
// @Tags         user
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Authentication header"
// @Success      200  {object}  model.UserInfoResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /user/info [get]
func (a *UserController) GetUserInfo(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	userIDAny, ok := c.Get(middleware.UserIDKey)
	if !ok {
		log.Warn("not JWT")
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}
	userID := userIDAny.(string)
	userGuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	user, serviceErr := a.service.GetInfo(userGuid, context.TODO())
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
		User: *user,
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
