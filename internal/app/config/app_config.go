package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest"
	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest/routes"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/supabase"

)

type StartUpConfig struct {
	DB  *gorm.DB
	App *gin.Engine
}

func StartUp(config *StartUpConfig) {
	// Third Party
	supabase := supabase.NewSupabaseClient()

	// Repository
	userRepository := repository.NewUserRepository(config.DB)
	eventRepository := repository.NewEventRepository(config.DB)
	cacheRepository := repository.NewCacheRepository()
	invoiceRepository := repository.NewInvoiceRepository(config.DB)

	// Service
	userService := service.NewUserService(userRepository, cacheRepository, supabase)
	eventService := service.NewEventService(eventRepository, supabase)
	paymentService := service.NewPaymentService(invoiceRepository, eventRepository)

	// Handler
	authHandler := rest.NewAuthHandler(userService)
	userHandler := rest.NewUserHandler(userService)
	eventHandler := rest.NewEventHandler(eventService)
	paymentHandler := rest.NewPaymentHandler(paymentService)

	routeSetting := routes.RouteConfig{
		App:          config.App,
		AuthHandler:  authHandler,
		UserHandler:  userHandler,
		EventHandler: eventHandler,
		PaymentHandler: paymentHandler,
	}
	routeSetting.Setup()
}
