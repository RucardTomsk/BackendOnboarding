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

// AddUser
// @Summary      Add User
// @Description  Add User
// @Tags         division
// @Accept       json
// @Produce      json
// @Param        payload body model.AddUserAndRoleRequest true "User request"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /division/addUser [post]
func (a *DivisionController) AddUser(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.AddUserAndRoleRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	divisionID, err := uuid.Parse(payload.DivisionID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	if err := a.service.AddUser(userID, divisionID, payload.Role, context.TODO()); err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
	})
}

// AddQuest
// @Summary      Add Quest
// @Description  Add Quest
// @Tags         quest
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateQuestRequest true "User request"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /quest/add [post]
func (a *DivisionController) AddQuest(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.CreateQuestRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	divisionID, err := uuid.Parse(payload.DivisionID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	if serviceErr := a.service.AddQuest(divisionID, &payload, context.TODO()); serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
	})
}

// AddStage
// @Summary      Add Stage
// @Description  Add Stage
// @Tags         quest
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateStageRequest true "User request"
// @Success      200  {object}  base.ResponseOK "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /quest/stage/add [post]
func (a *DivisionController) AddStage(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	var payload model.CreateStageRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		log.Warn("error parsing json:" + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	questID, err := uuid.Parse(payload.QuestID)
	if err != nil {
		log.Warn("error parsing uuid: " + err.Error())
		c.JSON(http.StatusBadRequest, api.GeneralParsingError(middleware.GetTrackingId(c)))
		return
	}

	if serviceErr := a.service.AddStage(questID, &payload, context.TODO()); serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, base.ResponseOK{
		Status:     http.StatusText(http.StatusOK),
		TrackingID: middleware.GetTrackingId(c),
	})
}

// GetAllQuest
// @Summary      Get All Quest
// @Description  Get All Quest
// @Tags         quest
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.GetAllQuestResponse "OK"
// @Failure      400  {object}  base.ResponseFailure "Bad request (client fault)"
// @Failure      500  {object}  base.ResponseFailure "Internal error (server fault)"
// @Router       /quest/all [get]
func (a DivisionController) GetAllQuest(c *gin.Context) {
	log := api.EnrichLogger(a.logger, c)

	guest, serviceErr := a.service.GetAllQuest()
	if serviceErr != nil {
		log.Warn("error occurred: " + serviceErr.Error())
		c.JSON(serviceErr.Code, api.ResponseFromServiceError(*serviceErr, middleware.GetTrackingId(c)))
		return
	}

	c.JSON(http.StatusOK, model.GetAllQuestResponse{
		ResponseOK: base.ResponseOK{
			Status:     http.StatusText(http.StatusOK),
			TrackingID: middleware.GetTrackingId(c),
		},
		Quests: guest,
	})
}
