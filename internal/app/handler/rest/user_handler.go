package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"

)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(us service.UserService) UserHandler {
	return UserHandler{
		userService: us,
	}
}

func (h *UserHandler) Current(ctx *gin.Context) {
	user := helper.GetUserLoginData(ctx)

	res, err := h.userService.UserCurrent(user)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}
