package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

func OrganizerRole() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {
		var res model.Response
		user := helper.GetUserLoginData(c)
		if !user.IsOrganizer {
			res.Error = true
			res.Message = "Only organizer can access this route"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		} else {
			c.Next()
		}
	})
}

func AdminRole() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {
		var res model.Response
		user := helper.GetUserLoginData(c)
		if !user.IsAdmin {
			res.Error = true
			res.Message = "Only admin can access this route"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		} else {
			c.Next()
		}
	})
}

func UserOnly() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {
		var res model.Response
		user := helper.GetUserLoginData(c)
		if user.IsAdmin || user.IsOrganizer {
			res.Error = true
			res.Message = "Admin and Organizer can't access this route"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		} else {
			c.Next()
		}
	})
}
