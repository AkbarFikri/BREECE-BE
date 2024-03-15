package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	ticketService  service.TicketService
}

func NewPaymentHandler(ps service.PaymentService, ts service.TicketService) PaymentHandler {
	return PaymentHandler{
		paymentService: ps,
		ticketService:  ts,
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

	res, err := h.paymentService.GenerateUrlAndToken(user, req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
	}

	helper.SuccessResponse(ctx, res)
}

func (h *PaymentHandler) Verify(ctx *gin.Context) {
	var notificationPayload map[string]interface{}

	err := ctx.ShouldBind(&notificationPayload)
	if err != nil {
		return
	}

	orderId, exists := notificationPayload["order_id"].(string)
	if !exists {
		return
	}

	success := h.paymentService.VerifyPayment(orderId)
	if !success {
		h.ticketService.FailurePayment(orderId)
	}

	h.ticketService.ConfirmedPayment(orderId)
}
