package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
)

type UserHandler struct {
	userService    service.UserService
	paymentService service.PaymentService
}

func NewUserHandler(us service.UserService, ps service.PaymentService) UserHandler {
	return UserHandler{
		userService:    us,
		paymentService: ps,
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

func (h *UserHandler) GetPaymentHistory(ctx *gin.Context) {
	user := helper.GetUserLoginData(ctx)

	res, err := h.paymentService.FetchPaymentHistory(user)

	if err != nil {
		helper.ErrorResponse(ctx, res)
	}

	helper.SuccessResponse(ctx, res)
}
