package appointment

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentControllerInterface interface {
	createAppointment(ctx *gin.Context)
	updateAppointmentStatus(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
}

func (c *AppointmentController) createAppointment(ctx *gin.Context) {
	var req createAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	res, err := c.service.CreateAppointment(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("create appointment successful", res))
}

func (c *AppointmentController) updateAppointmentStatus(ctx *gin.Context) {
	appointmentID := ctx.Param("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	var req updateAppointmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	// convert string to int64
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	err = c.service.UpdateAppointmentStatus(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("update appointment status successful", nil))
}

func (c *AppointmentController) getAppointmentByID(ctx *gin.Context) {
	appointmentID := ctx.Param("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	// convert string to int64
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetAppointmentByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get appointment by id successful", res))
}

func (c *AppointmentController) getAppointmentsByUser(ctx *gin.Context) {
	payload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetAppointmentsByUser(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get appointment successful", res))
}
