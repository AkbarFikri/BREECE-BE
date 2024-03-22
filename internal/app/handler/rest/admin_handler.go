package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/app/service"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/helper"
	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"
)

type AdminHandler struct {
	adminService service.AdminService
}

func NewAdminHandler(as service.AdminService) AdminHandler {
	return AdminHandler{
		adminService: as,
	}
}

func (h *AdminHandler) GetOrganizer(ctx *gin.Context) {
	id := ctx.Param("id")

	var res model.ServiceResponse
	var err error

	if id == "" {
		res, err = h.adminService.FetchOrganizer()
	} else {
		res, err = h.adminService.FetchOrganizerDetail(id)
	}

	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AdminHandler) VerifyOrganizer(ctx *gin.Context) {
	var req model.OrganizerVerifyRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid data provided",
		})
	}

	res, err := h.adminService.VerifyOrganizer(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}

func (h *AdminHandler) PostCategory(ctx *gin.Context) {
	var req model.CategoriesRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(ctx, model.ServiceResponse{
			Code:    http.StatusBadRequest,
			Error:   true,
			Message: "Invalid data provided",
		})
	}

	res, err := h.adminService.PostCategory(req)
	if err != nil {
		helper.ErrorResponse(ctx, res)
		return
	}

	helper.SuccessResponse(ctx, res)
}
