package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest"
	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest/middleware"

)

type RouteConfig struct {
	App            *gin.Engine
	UserHandler    rest.UserHandler
	AuthHandler    rest.AuthHandler
	EventHandler   rest.EventHandler
	PaymentHandler rest.PaymentHandler
	AdminHandler   rest.AdminHandler
}

func (c *RouteConfig) Setup() {
	c.ServeRoute()
}

func (c *RouteConfig) ServeRoute() {
	c.App.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	c.App.Use(gin.Logger())
	c.App.Use(gin.Recovery())
	c.App.Use(middleware.CORSMiddleware())
	v1 := c.App.Group("/api/v1")
	v1.StaticFS("/docs", http.Dir("api/dist"))
	v1.StaticFS("/doc", http.Dir("api"))
	c.AuthRoute(v1)
	c.UserRoute(v1)
	c.EventRoute(v1)
	c.PaymentRoute(v1)
}

func (c *RouteConfig) AuthRoute(r *gin.RouterGroup) {
	authEnds := r.Group("/auth")
	authEnds.GET("/check", c.AuthHandler.HealthCheck)
	authEnds.POST("/register", c.AuthHandler.RegisterUser)
	authEnds.POST("/register/organizer", c.AuthHandler.RegisterOrganizer)
	authEnds.POST("/register/admin", middleware.APIKEY(), c.AuthHandler.RegisterAdmin)
	authEnds.POST("/login", c.AuthHandler.LoginUser)
	authEnds.POST("/login/organizer", c.AuthHandler.LoginOrganizer)
	authEnds.POST("/login/admin", c.AuthHandler.LoginAdmin)
	authEnds.POST("/otp", c.AuthHandler.VerifyOTP)
	authEnds.POST("/profile", c.AuthHandler.VerifyProfile)
}

func (c *RouteConfig) UserRoute(r *gin.RouterGroup) {
	userEnds := r.Group("/user")
	userEnds.Use(middleware.JwtUser())
	userEnds.GET("/current", c.UserHandler.Current)
	userEnds.GET("/payment", c.UserHandler.GetPaymentHistory)
	userEnds.GET("/event", c.UserHandler.GetTicketHisoty)
}

func (c *RouteConfig) EventRoute(r *gin.RouterGroup) {
	eventEnds := r.Group("/event")
	eventEnds.Use(middleware.JwtUser())
	eventEnds.POST("/", middleware.OrganizerRole(), c.EventHandler.PostEvent)
	eventEnds.GET("/", c.EventHandler.GetEvent)
	eventEnds.GET("/category", c.EventHandler.GetEventCategory)
	eventEnds.GET("/:id", c.EventHandler.GetEventDetails)
	eventEnds.GET("/:id/participant", middleware.OrganizerRole(), c.EventHandler.GetEventParticipant)
}

func (c *RouteConfig) PaymentRoute(r *gin.RouterGroup) {
	paymentEnds := r.Group("/payment")
	paymentEnds.POST("/checkout", middleware.JwtUser(), middleware.UserOnly(), c.PaymentHandler.Checkout)
	paymentEnds.POST("/verify", c.PaymentHandler.Verify)
}

func (c *RouteConfig) AdminRoute(r *gin.RouterGroup) {
	adminEnds := r.Group("/admin")
	adminEnds.Use(middleware.JwtUser())
	adminEnds.Use(middleware.AdminRole())
	adminEnds.GET("/organizer", c.AdminHandler.GetOrganizer)
	adminEnds.GET("/organizer/:id", c.AdminHandler.GetOrganizer)
	adminEnds.PATCH("/organizer/verify", c.AdminHandler.VerifyOrganizer)
}
