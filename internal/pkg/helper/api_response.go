package helper

import (
	"github.com/gin-gonic/gin"

	"github.com/AkbarFikri/BREECE-BE/internal/pkg/model"

)

func SuccessResponse(c *gin.Context, data model.ServiceResponse) {
	res := &model.Response{
		Error:   data.Error,
		Message: data.Message,
		Data:    data.Data,
	}

	c.JSON(data.Code, res)
	return
}

func ErrorResponse(c *gin.Context, data model.ServiceResponse) {
	res := &model.Response{
		Error:   data.Error,
		Message: data.Message,
		Data:    data.Data,
	}

	c.AbortWithStatusJSON(data.Code, res)
	return
}
