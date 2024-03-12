package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(es service.EventService) EventHandler {
	return EventHandler{
		eventService: es,
	}
}

func (h *EventHandler) PostEvent(ctx *gin.Context) {
	user := helper.GetUserLoginData(ctx)

	var req model.EventRequest

	if err := ctx.ShouldBind(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
			Data:    err.Error(),
		})
		return
	}

	data, err := h.eventService.PostEvent(user, req)
	if err != nil {
		helper.ErrorResponse(ctx, data)
	}

	helper.SuccessResponse(ctx, data)
}
