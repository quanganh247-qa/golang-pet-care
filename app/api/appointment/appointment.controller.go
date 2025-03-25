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
	confirmAppointment(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
	getAppointmentsByDoctor(ctx *gin.Context)
	//time slot
	getAvailableTimeSlots(ctx *gin.Context)
}

func (c *AppointmentController) createAppointment(ctx *gin.Context) {
	var req createAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := c.service.CreateAppointment(ctx, req, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("create appointment successful", res))
}

func (c *AppointmentController) confirmAppointment(ctx *gin.Context) {
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
	err = c.service.ConfirmPayment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("confirm appointment successful", nil))
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

func (c *AppointmentController) getAppointmentsByDoctor(ctx *gin.Context) {
	doctorID := ctx.Param("doctor_id")
	if doctorID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	// convert string to int64
	id, err := strconv.ParseInt(doctorID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetAppointmentsByDoctor(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get appointment successful", res))
}

func (c *AppointmentController) getAvailableTimeSlots(ctx *gin.Context) {

	doctorID := ctx.Param("doctor_id")
	if doctorID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	date := ctx.Query("date")

	id, err := strconv.ParseInt(doctorID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetAvailableTimeSlots(ctx, id, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get time slot successful", res))
}
