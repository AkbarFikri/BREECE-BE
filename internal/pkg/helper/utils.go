package helper

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

func GenerateRandomString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRandomInt(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var letterRunes = []rune("0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetUserLoginData(c *gin.Context) model.UserTokenData {
	getUser, _ := c.Get("user")
	user := getUser.(model.UserTokenData)

	return user
}
