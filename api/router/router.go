package router

import (
	"github.com/RucardTomsk/BackendOnboarding/api/controller"
	"github.com/RucardTomsk/BackendOnboarding/cmd/config"
	"github.com/RucardTomsk/BackendOnboarding/internal/api/middleware"
	"github.com/RucardTomsk/BackendOnboarding/internal/common"
	"github.com/RucardTomsk/BackendOnboarding/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	config config.Config
}

func NewRouter(config config.Config) *Router {
	return &Router{
		config: config,
	}
}

func (h *Router) InitRoutes(
	logger *zap.Logger,
	userService *service.UserService,
	controllerContainer *controller.Container) *gin.Engine {

	gin.SetMode(h.config.Server.GinMode)

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(middleware.SetTracingContext(*logger))
	router.Use(middleware.SetRequestLogging(*logger))
	router.Use(middleware.SetOperationName(h.config.Server, *logger))
	router.Use(middleware.SetAccessControl(h.config.Server, *logger))
	router.Use(cors.New(common.DefaultCorsConfig()))

	router.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	baseRouter := router.Group("/api")
	v1 := baseRouter.Group("/v1")
	{
		v1.GET("test", middleware.AuthorizationCheck(userService, *logger), controllerContainer.UserController.Test)
	}
	user := v1.Group("user")
	{
		user.POST("register", controllerContainer.UserController.RegistrationUser)
		user.POST("login", controllerContainer.UserController.LoginUser)
		user.GET("info", middleware.AuthorizationCheck(userService, *logger), controllerContainer.UserController.GetUserInfo)
		user.GET("all", controllerContainer.UserController.GetUsers)
		user.POST("about/update", middleware.AuthorizationCheck(userService, *logger), controllerContainer.UserController.UpdateAbout)
		user.GET("quest", middleware.AuthorizationCheck(userService, *logger), controllerContainer.UserController.GetUserQuest)
	}

	division := v1.Group("division")
	{
		division.POST("create", controllerContainer.DivisionController.CreateDivision)
		division.GET("all", controllerContainer.DivisionController.GetDivisions)
		division.POST("addUser", controllerContainer.DivisionController.AddUser)
	}

	quest := v1.Group("quest")
	{
		quest.POST("add", controllerContainer.DivisionController.AddQuest)
		quest.POST("stage/add", controllerContainer.DivisionController.AddStage)
		quest.GET("all", controllerContainer.DivisionController.GetAllQuest)
	}

	roles := v1.Group("roles")
	{
		roles.GET("all", controllerContainer.RolesController.GetRoles)
	}

	return router
}
