package appointment

import (
	"net/http"
<<<<<<< HEAD
	"strconv"
=======
>>>>>>> 323513c (appointment api)

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentControllerInterface interface {
	createAppointment(ctx *gin.Context)
<<<<<<< HEAD
	updateAppointmentStatus(ctx *gin.Context)
	getAppointmentsOfDoctor(ctx *gin.Context)
=======
>>>>>>> 323513c (appointment api)
}

func (c *AppointmentController) createAppointment(ctx *gin.Context) {
	var req createAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
<<<<<<< HEAD
=======

>>>>>>> 323513c (appointment api)
	res, err := c.service.CreateAppointment(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, util.SuccessResponse("create appointment successful", res))
<<<<<<< HEAD
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

func (c *AppointmentController) getAppointmentsOfDoctor(ctx *gin.Context) {

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
	res, err := c.service.GetAppointmentsOfDoctorService(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get appointments of doctor successful", res))
=======
>>>>>>> 323513c (appointment api)
}
