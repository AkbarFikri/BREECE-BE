package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type PaymentHandler struct {
	PaymentService service.PaymentService
}

func NewPaymentHandler(ps service.PaymentService) PaymentHandler {
	return PaymentHandler{
		PaymentService: ps,
	}
}

func (h *PaymentHandler) Checkout(ctx *gin.Context) {
	user := helper.GetUserLoginData(ctx)

	var req model.PaymentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid data provided",
		})
		return
	}

	// ctx.JSON(200, gin.H{
	// 	"user": user,
	// 	"data": req,
	// })

	res, err := h.PaymentService.GenerateUrlAndToken(user, req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
	}

	helper.SuccessResponse(ctx, res)
}
