package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest"
	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest/routes"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
)

type StartUpConfig struct {
	DB  *gorm.DB
	App *gin.Engine
}

func StartUp(config *StartUpConfig) {
	// Repository
	userRepository := repository.NewUserRepository(config.DB)
	cacheRepository := repository.NewCacheRepository()

	// Service
	userService := service.NewUserService(userRepository, cacheRepository)

	// Handler
	authHandler := rest.NewAuthHandler(userService)

	routeSetting := routes.RouteConfig{
		App:         config.App,
		AuthHandler: authHandler,
	}
	routeSetting.Setup()
}
