package service_type

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ServiceTypeControllerInterface interface {
	CreateServiceType(ctx *gin.Context)
}

func (controller *ServiceTypeController) CreateServiceType(ctx *gin.Context) {
	var req ServiceType
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	response, err := controller.service.createServiceTypeService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", response))

}
