package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

type EventHandler struct {
	eventService  service.EventService
	ticketService service.TicketService
}

func NewEventHandler(es service.EventService, ts service.TicketService) EventHandler {
	return EventHandler{
		eventService:  es,
		ticketService: ts,
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

	res, err := h.eventService.PostEvent(user, req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *EventHandler) GetEvent(ctx *gin.Context) {
	user := helper.GetUserLoginData(ctx)

	var params model.FilterParam
	if err := ctx.ShouldBindQuery(&params); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusInternalServerError,
			Error:   true,
			Message: "Something went wrong",
		})
	}

	res, err := h.eventService.FetchEvent(user, params)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *EventHandler) GetEventDetails(ctx *gin.Context) {
	id := ctx.Param("id")

	if id == "" {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
		})
		return
	}

	res, err := h.eventService.FetchEventDetails(id)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *EventHandler) GetEventCategory(ctx *gin.Context) {
	res, err := h.eventService.FetchCategory()
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *EventHandler) GetEventParticipant(ctx *gin.Context) {
	eventId := ctx.Param("id")
	user := helper.GetUserLoginData(ctx)

	if eventId == "" {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Bad request, data provided is invalid",
		})
		return
	}

	res, err := h.ticketService.FetchParticipantTicket(user, eventId)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)

}
