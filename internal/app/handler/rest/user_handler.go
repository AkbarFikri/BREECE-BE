package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
)

type UserHandler struct {
	userService service.UserService
}

func (h *UserHandler) current(ctx *gin.Context) {

}
