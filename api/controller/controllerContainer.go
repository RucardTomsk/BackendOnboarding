package controller

import (
	"github.com/RucardTomsk/BackendOnboarding/service"
	"go.uber.org/zap"
)

type Container struct {
	UserController     *UserController
	DivisionController *DivisionController
	RolesController    *RolesController
}

func NewControllerContainer(
	logger *zap.Logger,
	userService *service.UserService,
	divisionService *service.DivisionService) *Container {
	return &Container{
		UserController:     NewUserController(logger, userService),
		DivisionController: NewDivisionController(logger, divisionService),
		RolesController:    NewRolesController(logger),
	}
}
