package controller

import (
	"github.com/RucardTomsk/BackendOnboarding/service"
	"go.uber.org/zap"
)

type Container struct {
	UserController *UserController
}

func NewControllerContainer(
	logger *zap.Logger,
	userService *service.UserService) *Container {
	return &Container{
		UserController: NewUserController(logger, userService),
	}
}
