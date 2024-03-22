package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

func APIKEY() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {
		var res model.Response

		reqKey := c.GetHeader("KEY")

		if reqKey == "" || reqKey != os.Getenv("API_KEY") {
			res.Error = true
			res.Message = "Unauthorized"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		} else {
			c.Next()
		}
	})
}
