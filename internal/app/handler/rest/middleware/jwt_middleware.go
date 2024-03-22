package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

func JwtUser() gin.HandlerFunc {

	return gin.HandlerFunc(func(c *gin.Context) {
		var res model.Response

		if c.GetHeader("Authorization") == "" {
			res.Error = true
			res.Message = "Authorization is required for this endpoint"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusForbidden, res)
		}

		if !strings.Contains(c.GetHeader("Authorization"), "Bearer") {
			res.Error = true
			res.Message = "accessToken invalid or expired"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusForbidden, res)
		}

		token, err := helper.VerifyTokenHeader(c)

		if err != nil {
			res.Error = true
			res.Message = "accessToken invalid or expired"
			res.Data = nil
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		} else {
			claims := token.Claims.(jwt.MapClaims)
			user := model.UserTokenData{
				ID:                claims["id"].(string),
				Email:             claims["email"].(string),
				IsAdmin:           claims["isAdmin"].(bool),
				IsEmailVerified:   claims["isEmailVerified"].(bool),
				IsProfileVerified: claims["isProfileVerified"].(bool),
				IsOrganizer:       claims["isOrganizer"].(bool),
				IsBrawijaya:       claims["isBrawijaya"].(bool),
			}
			c.Set("user", user)
			c.Next()
		}
	})

}
