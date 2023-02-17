package main

import (
	"context"
	"fmt"
	"github.com/RucardTomsk/BackendOnboarding/api/controller"
	"github.com/RucardTomsk/BackendOnboarding/api/router"
	"github.com/RucardTomsk/BackendOnboarding/cmd/config"
	"github.com/RucardTomsk/BackendOnboarding/cmd/docs"
	"github.com/RucardTomsk/BackendOnboarding/internal/common"
	"github.com/RucardTomsk/BackendOnboarding/internal/server"
	"github.com/RucardTomsk/BackendOnboarding/internal/telemetry/log"
	"github.com/RucardTomsk/BackendOnboarding/service"
	postgresStorage "github.com/RucardTomsk/BackendOnboarding/storage/dao/postgres"
	"github.com/RucardTomsk/BackendOnboarding/storage/migration"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.NewLogger()

	appCli := common.InitAppCli()
	if err := appCli.Run(os.Args); err != nil {
		logger.Fatal(err.Error())
	}

	// read config
	var cfg config.Config
	if err := viper.MergeInConfig(); err != nil {
		logger.Fatal(fmt.Sprintf("error reading config file: %v", err))
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		logger.Fatal(fmt.Sprintf("unable to decode into struct: %v", err))
	}

	// configure swagger
	swaggerConfig := common.NewSwaggerConfig("Task API", "TBD", "unreleased")

	docs.SwaggerInfo.Title = swaggerConfig.Title
	docs.SwaggerInfo.Description = swaggerConfig.Description
	docs.SwaggerInfo.Version = swaggerConfig.Version
	docs.SwaggerInfo.Host = swaggerConfig.Host
	docs.SwaggerInfo.BasePath = swaggerConfig.BasePath
	docs.SwaggerInfo.Schemes = swaggerConfig.Schemes

	// init connections
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Name, cfg.Postgres.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal(fmt.Sprintf("can't connect to database: %v", err))
	}

	logger.Info(fmt.Sprintf("successfully connected to database %s on %s:%d as %s",
		cfg.Postgres.Name, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User))

	if err := migration.Migrate(db); err != nil {
		logger.Fatal(fmt.Sprintf("failed to migrate database: %v", err))
	}
	logger.Info("database migrated successfully")

	// init storage
	userStorage := postgresStorage.NewUserStorage(db)

	// init services
	userService := service.NewUserService(userStorage)

	// init controllers
	controllers := controller.NewControllerContainer(
		logger,
		userService)

	// init data processing

	// init server
	handler := router.NewRouter(cfg)
	srv := new(server.Server)

	go func() {
		if err := srv.Run(cfg.Server.Host, cfg.Server.Port, handler.InitRoutes(
			logger,
			userService,
			controllers)); err != nil {
			logger.Error(fmt.Sprintf("error accured while running http server: %s", err.Error()))
		}
	}()

	logger.Info(fmt.Sprintf("listening on %s:%s", cfg.Server.Host, cfg.Server.Port))

	// handle signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Info("shutting down gracefully...")
	defer func() { logger.Info("shutdown complete") }()

	// perform shutdown
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}
}
