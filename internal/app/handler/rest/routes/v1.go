package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest"
	"github.com/AkbarFikri/BREECE-BE/internal/app/handler/rest/middleware"
)

type RouteConfig struct {
	App         *gin.Engine
	UserHandler rest.UserHandler
	AuthHandler rest.AuthHandler
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

}

func (c *RouteConfig) AuthRoute(r *gin.RouterGroup) {
	authEnds := r.Group("/auth")
	authEnds.GET("/check", c.AuthHandler.HealthCheck)
	authEnds.POST("/register", c.AuthHandler.Register)
	authEnds.POST("/login", c.AuthHandler.Login)
	authEnds.POST("/otp", c.AuthHandler.VerifyOTP)
	authEnds.POST("/profile", c.AuthHandler.VerifyProfile)
}

func (c *RouteConfig) UserRoute(r *gin.RouterGroup) {
	userEnds := r.Group("/user")
	userEnds.Use(middleware.JwtUser())
	userEnds.GET("/current", c.UserHandler.Current)
}
