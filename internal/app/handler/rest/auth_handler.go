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

func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
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

	req.IsOrganizer = false

	res, err := h.userService.Register(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) RegisterOrganizer(ctx *gin.Context) {
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

	req.IsOrganizer = true

	res, err := h.userService.Register(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) RegisterAdmin(ctx *gin.Context) {
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

	res, err := h.userService.RegisterAdmin(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) LoginUser(ctx *gin.Context) {
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

	res, err := h.userService.LoginUser(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)

}

func (h *AuthHandler) LoginOrganizer(ctx *gin.Context) {
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

	res, err := h.userService.LoginOrganizer(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)

}

func (h *AuthHandler) LoginAdmin(ctx *gin.Context) {
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

	res, err := h.userService.LoginAdmin(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) VerifyOTP(ctx *gin.Context) {
	var req model.OtpUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
			Data:    nil,
		})
		return
	}

	res, err := h.userService.VerifyOTP(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AuthHandler) VerifyProfile(ctx *gin.Context) {
	var req model.ProfileUserRequest

	if err := ctx.ShouldBind(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
			Data:    nil,
		})
		return
	}

	res, err := h.userService.VerifyProfile(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}
