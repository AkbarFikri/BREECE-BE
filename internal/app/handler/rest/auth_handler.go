package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(us service.UserService) AuthHandler {
	return AuthHandler{
		userService: us,
	}
}

func (h *AuthHandler) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"Status": "Ready To Code!"})
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req model.CreateUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
			Data:    nil,
		})
		return
	}

	res, err := h.userService.Register(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req model.LoginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
			Data:    nil,
		})
		return
	}

	res, err := h.userService.Login(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)

}
