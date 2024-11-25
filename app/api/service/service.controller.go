package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type ServiceControllerInterface interface {
	CreateServiceController(ctx *gin.Context)
	DeleteService(ctx *gin.Context)
	GetAllServices(ctx *gin.Context)
	UpdateService(ctx *gin.Context)
	GetServiceByID(ctx *gin.Context)
	getAllServices(ctx *gin.Context)
}

func (controller *ServiceController) CreateServiceController(ctx *gin.Context) {
	var req CreateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	response, err := controller.service.CreateService(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("Success", response))

}

// func (controller *ServiceController) DeleteService(ctx *gin.Context) {
// 	serviceID, err := strconv.ParseInt(ctx.Param("serviceid"), 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	err = controller.service.deleteServiceService(ctx, serviceID)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
// }

func (controller *ServiceController) GetAllServices(ctx *gin.Context) {

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
<<<<<<< HEAD
	services, err := controller.service.GetAllServices(ctx, pagination)
=======
	services, err := controller.service.getAllServicesService(ctx, pagination)
>>>>>>> c73e2dc (pagination function)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, services)
}

func (controller *ServiceController) GetServiceByID(ctx *gin.Context) {

	serviceidStr := ctx.Param("id")

	serviceid, err := strconv.ParseInt(serviceidStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.GetServiceByID(ctx, serviceid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

func (controller *ServiceController) UpdateService(ctx *gin.Context) {
	serviceid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}
	var req UpdateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := controller.service.UpdateService(ctx, serviceid, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", res))
}

// Delete service
func (controller *ServiceController) DeleteService(ctx *gin.Context) {
	serviceID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = controller.service.DeleteService(ctx, serviceID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Success", nil))
}

func (controller *ServiceController) getAllServices(ctx *gin.Context) {
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	services, err := controller.service.getAllServices(ctx, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, services)
}
