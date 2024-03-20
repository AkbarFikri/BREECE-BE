package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest"
	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest/routes"
	"github.com/AkbarFikri/BREECE-BE/internal/app/repository"
	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/mailer"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/supabase"

)

type StartUpConfig struct {
	DB  *gorm.DB
	App *gin.Engine
}

func StartUp(config *StartUpConfig) {
	// Third Party
	supabase := supabase.NewSupabaseClient()
	mailer := mailer.NewMailer()

	// Repository
	userRepository := repository.NewUserRepository(config.DB)
	eventRepository := repository.NewEventRepository(config.DB)
	cacheRepository := repository.NewCacheRepository()
	invoiceRepository := repository.NewInvoiceRepository(config.DB)
	ticketRepository := repository.NewTicketRepository(config.DB)

	// Service
	userService := service.NewUserService(userRepository, cacheRepository, supabase, mailer)
	eventService := service.NewEventService(eventRepository, supabase)
	paymentService := service.NewPaymentService(invoiceRepository, eventRepository)
	ticketService := service.NewTicketService(eventRepository, invoiceRepository, ticketRepository, userRepository, mailer)
	adminService := service.NewAdminService(userRepository)

	// Handler
	authHandler := rest.NewAuthHandler(userService)
	userHandler := rest.NewUserHandler(userService)
	eventHandler := rest.NewEventHandler(eventService)
	paymentHandler := rest.NewPaymentHandler(paymentService, ticketService)
	adminHandler := rest.NewAdminHandler(adminService)

	routeSetting := routes.RouteConfig{
		App:            config.App,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		EventHandler:   eventHandler,
		PaymentHandler: paymentHandler,
		AdminHandler:   adminHandler,
	}
	routeSetting.Setup()
}
