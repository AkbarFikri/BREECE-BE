package helper

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func SignJWT(Data map[string]interface{}, ExpiredAt time.Duration) (string, error) {
	expiredAt := time.Now().Add(time.Duration(time.Hour) * ExpiredAt).Unix()

	JWTSecretKey := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{}
	claims["exp"] = expiredAt
	claims["authorization"] = true

	for i, v := range Data {
		claims[i] = v
	}

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := to.SignedString([]byte(JWTSecretKey))

	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

func VerifyTokenHeader(c *gin.Context) (*jwt.Token, error) {
	header := c.GetHeader("Authorization")
	accessToken := strings.SplitAfter(header, "Bearer")[1]
	JWTSecretKey := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(strings.Trim(accessToken, " "), func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
