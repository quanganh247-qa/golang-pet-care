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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	confirmAppointment(ctx *gin.Context)
	checkinAppointment(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
	getAppointmentsByDoctor(ctx *gin.Context)
	getAllAppointments(ctx *gin.Context)
	getAllAppointmentsByDate(ctx *gin.Context)
	updateAppointment(ctx *gin.Context)
	getQueue(ctx *gin.Context)
	updateQueueItemStatus(ctx *gin.Context)
	getHistoryAppointmentsByPetID(ctx *gin.Context)

	//time slot
	getAvailableTimeSlots(ctx *gin.Context)

	// SOAP
	createSOAP(ctx *gin.Context)
	updateSOAP(ctx *gin.Context)
	getSOAPByAppointmentID(ctx *gin.Context)
=======
	updateAppointmentStatus(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
	getAppointmentByID(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)
=======
	getAppointmentsByPetOfUser(ctx *gin.Context)
>>>>>>> e30b070 (Get list appoinment by user)
=======
	getAppointmentsByUser(ctx *gin.Context)
>>>>>>> 685da65 (latest update)
=======
	confirmAppointment(ctx *gin.Context)
	checkinAppointment(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
	getAppointmentsByDoctor(ctx *gin.Context)
	getAllAppointments(ctx *gin.Context)
	getAllAppointmentsByDate(ctx *gin.Context)
	updateAppointment(ctx *gin.Context)
	getQueue(ctx *gin.Context)
	updateQueueItemStatus(ctx *gin.Context)
	getHistoryAppointmentsByPetID(ctx *gin.Context)

	//time slot
	getAvailableTimeSlots(ctx *gin.Context)
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======

	// SOAP
	createSOAP(ctx *gin.Context)
	updateSOAP(ctx *gin.Context)
<<<<<<< HEAD
>>>>>>> e859654 (Elastic search)
=======
	getSOAPByAppointmentID(ctx *gin.Context)
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
	getAppointmentsOfDoctor(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)
=======
	getAppointmentsByPetOfUser(ctx *gin.Context)
>>>>>>> e30b070 (Get list appoinment by user)
=======
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
>>>>>>> 685da65 (latest update)
=======
	confirmAppointment(ctx *gin.Context)
	getAppointmentByID(ctx *gin.Context)
	getAppointmentsByUser(ctx *gin.Context)
	getAppointmentsByDoctor(ctx *gin.Context)
	getAllAppointments(ctx *gin.Context)
	//time slot
	getAvailableTimeSlots(ctx *gin.Context)
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======

	// SOAP
	createSOAP(ctx *gin.Context)
	updateSOAP(ctx *gin.Context)
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	appointmentID := ctx.Param("id")
=======
=======
>>>>>>> b393bb9 (add service and add permission)
	appointmentID := ctx.Param("appointment_id")
>>>>>>> b393bb9 (add service and add permission)
=======
	appointmentID := ctx.Param("id")
>>>>>>> e859654 (Elastic search)
=======
	appointmentID := ctx.Param("id")
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD

=======
>>>>>>> b393bb9 (add service and add permission)
=======

>>>>>>> e859654 (Elastic search)
=======
>>>>>>> b393bb9 (add service and add permission)
=======

>>>>>>> e859654 (Elastic search)
	err = c.service.ConfirmPayment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("confirm appointment successful", nil))
}

<<<<<<< HEAD
<<<<<<< HEAD
func (c *AppointmentController) checkinAppointment(ctx *gin.Context) {
	roomID, err := strconv.ParseInt(ctx.Query("room_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// convert string to int64
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	priority := ctx.Query("priority")

	err = c.service.CheckInAppoinment(ctx, id, roomID, priority)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("checkin appointment successful", nil))
}

func (c *AppointmentController) getAppointmentByID(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
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

<<<<<<< HEAD
func (c *AppointmentController) getAppointmentsByUser(ctx *gin.Context) {
	payload, err := middleware.GetAuthorizationPayload(ctx)
=======
func (c *AppointmentController) checkinAppointment(ctx *gin.Context) {
	roomID, err := strconv.ParseInt(ctx.Query("room_id"), 10, 64)
>>>>>>> 4ccd381 (Update appointment flow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
<<<<<<< HEAD
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
=======
	// convert string to int64
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
>>>>>>> 4ccd381 (Update appointment flow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
<<<<<<< HEAD
	res, err := c.service.GetAppointmentsByDoctor(ctx, id)
=======

	priority := ctx.Query("priority")

	err = c.service.CheckInAppoinment(ctx, id, roomID, priority)
>>>>>>> 4ccd381 (Update appointment flow)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
<<<<<<< HEAD
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

func (c *AppointmentController) getAllAppointments(ctx *gin.Context) {

	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	option := ctx.Query("option")

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetAllAppointments(ctx, date, option, pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get all appointments successful", res))
}

func (c *AppointmentController) getAllAppointmentsByDate(ctx *gin.Context) {

	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetAllAppointmentsByDate(ctx, pagination, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get all appointments by date successful", res))
}
func (c *AppointmentController) updateAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var req updateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	err = c.service.UpdateAppointmentService(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("update appointment successful", nil))
}

func (c *AppointmentController) createSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap CreateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	soapResponse, err := c.service.CreateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}

func (c *AppointmentController) updateSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap UpdateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	soapResponse, err := c.service.UpdateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}

func (c *AppointmentController) getQueue(ctx *gin.Context) {

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	queueItems, err := c.service.GetQueueService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, queueItems)
}

func (c *AppointmentController) updateQueueItemStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateQueueItemStatusService(ctx, id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (c *AppointmentController) getHistoryAppointmentsByPetID(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	if petID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	historyAppointments, err := c.service.GetHistoryAppointmentsByPetID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, historyAppointments)
}

func (c *AppointmentController) getSOAPByAppointmentID(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	soap, err := c.service.GetSOAPByAppointmentID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, soap)
}

=======
>>>>>>> 685da65 (latest update)
=======
	ctx.JSON(http.StatusOK, util.SuccessResponse("checkin appointment successful", nil))
}

>>>>>>> 4ccd381 (Update appointment flow)
func (c *AppointmentController) getAppointmentByID(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
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

func (c *AppointmentController) getAllAppointments(ctx *gin.Context) {

	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	option := ctx.Query("option")

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetAllAppointments(ctx, date, option, pagination)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get all appointments successful", res))
}

func (c *AppointmentController) getAllAppointmentsByDate(ctx *gin.Context) {

	date := ctx.Query("date")
	if date == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetAllAppointmentsByDate(ctx, pagination, date)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get all appointments by date successful", res))
}
func (c *AppointmentController) updateAppointment(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var req updateAppointmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	err = c.service.UpdateAppointmentService(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("update appointment successful", nil))
}

func (c *AppointmentController) createSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap CreateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	soapResponse, err := c.service.CreateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}

func (c *AppointmentController) updateSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap UpdateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	soapResponse, err := c.service.UpdateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}

func (c *AppointmentController) getQueue(ctx *gin.Context) {

	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	queueItems, err := c.service.GetQueueService(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, queueItems)
}

func (c *AppointmentController) updateQueueItemStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.service.UpdateQueueItemStatusService(ctx, id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (c *AppointmentController) getHistoryAppointmentsByPetID(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	if petID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	historyAppointments, err := c.service.GetHistoryAppointmentsByPetID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, historyAppointments)
}

func (c *AppointmentController) getSOAPByAppointmentID(ctx *gin.Context) {
	appointmentID := ctx.Param("id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	soap, err := c.service.GetSOAPByAppointmentID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, soap)
}

=======
>>>>>>> 685da65 (latest update)
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

func (c *AppointmentController) getAllAppointments(ctx *gin.Context) {
	res, err := c.service.GetAllAppointments(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("get all appointments successful", res))
}

func (c *AppointmentController) createSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap CreateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	soapResponse, err := c.service.CreateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}

func (c *AppointmentController) updateSOAP(ctx *gin.Context) {
	appointmentID := ctx.Param("appointment_id")
	if appointmentID == "" {
		ctx.JSON(http.StatusBadRequest, nil)
		return
	}
	id, err := strconv.ParseInt(appointmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	var soap UpdateSOAPRequest
	if err := ctx.ShouldBindJSON(&soap); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}
	soapResponse, err := c.service.UpdateSOAPService(ctx, soap, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP", soapResponse))
}
