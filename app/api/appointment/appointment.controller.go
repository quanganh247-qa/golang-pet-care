package appointment

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type AppointmentControllerInterface interface {
	createAppointment(ctx *gin.Context)
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
	getCompletedAppointments(ctx *gin.Context)
	//time slot
	getAvailableTimeSlots(ctx *gin.Context)

	// SOAP
	createSOAP(ctx *gin.Context)
	updateSOAP(ctx *gin.Context)
	getSOAPByAppointmentID(ctx *gin.Context)

	// statistic
	GetAppointmentDistribution(ctx *gin.Context)

	CreateWalkInAppointment(c *gin.Context)
	GetAppointmentByState(ctx *gin.Context)
	HandleWebSocket(ctx *gin.Context)

	// Long polling
	waitForNotifications(ctx *gin.Context)
	getNotifications(ctx *gin.Context)
	getMissedNotifications(ctx *gin.Context)

	// User role registration
	registerUserRole(ctx *gin.Context)

	// Notifications from DB
	getNotificationsFromDB(ctx *gin.Context)
	markNotificationAsRead(ctx *gin.Context)
	markAllNotificationsAsRead(ctx *gin.Context)
}

// createAppointment creates a new appointment
// @Summary Create a new appointment
// @Description Create a new appointment for a user
// @Accept json
// @Produce json
// @Param appointment body createAppointmentRequest true "Appointment details"
// @Success 200 {object} util.SuccessResponse{data=createAppointmentResponse} "Appointment created successfully"
// @Failure 400 {object} util.ErrorResponse "Invalid request"
// @Failure 500 {object} util.ErrorResponse "Internal server error"
// @Router /appointment [post]
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

	err = c.service.ConfirmPayment(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("confirm appointment successful", nil))
}

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

// getAppointmentByID retrieves an appointment by its ID
// @Summary Get an appointment by ID
// @Description Retrieve an appointment by its unique identifier
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID"
// @Success 200 {object} util.SuccessResponse{data=Appointment} "Appointment retrieved successfully"
// @Failure 400 {object} util.ErrorResponse "Invalid request"
// @Failure 500 {object} util.ErrorResponse "Internal server error"
// @Router /appointment/{id} [get]
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
	var date string
	option := ctx.Query("option")
	if option == "false" {

		date = ctx.Query("date")
		if date == "" {
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

	}

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

func (c *AppointmentController) GetAppointmentDistribution(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	distributions, err := c.service.GetAppointmentDistributionByService(ctx, startDateStr, endDateStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    distributions,
	})
}

func (ac *AppointmentController) CreateWalkInAppointment(c *gin.Context) {
	var req createWalkInAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	// Validate the request
	if req.DoctorID == 0 {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("doctor ID is required")))
		return
	}

	if req.ServiceID == 0 {
		c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("service ID is required")))
		return
	}

	// Case 1: New user with new pet
	if req.Owner != nil {
		if req.Owner.OwnerName == "" || req.Owner.OwnerPhone == "" {
			c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("owner name and phone are required for new users")))
			return
		}

		if req.Pet == nil {
			c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("pet information is required for new users")))
			return
		}

		if req.Pet.Name == "" || req.Pet.Species == "" {
			c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("pet name and species are required")))
			return
		}
	} else {
		// Case 2: Existing pet
		if req.PetID == 0 {
			c.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("pet ID is required for existing users")))
			return
		}
	}

	// Create walk-in appointment
	response, err := ac.service.CreateWalkInAppointment(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// HandleWebSocket handles WebSocket connections
func (c *AppointmentController) HandleWebSocket(ctx *gin.Context) {
	// Get user information for client ID
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Set client ID in query parameters
	clientID := fmt.Sprintf("user_%s", authPayload.Username)
	ctx.Request.Header.Set("X-Client-ID", clientID)

	// Handle WebSocket connection
	c.service.HandleWebSocket(ctx)
}

func (c *AppointmentController) GetAppointmentByState(ctx *gin.Context) {
	// Get the state from the query parameters
	state := ctx.Query("state")
	if state == "" {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(fmt.Errorf("state is required")))
		return
	}

	// Get the appointments by state
	appointments, err := c.service.GetAppointmentByState(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}

// WaitForNotifications thực hiện Long polling để chờ thông báo mới
// Phương thức này sẽ giữ kết nối mở trong thời gian chờ quy định (30s)
func (c *AppointmentController) waitForNotifications(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Chuyển các thông báo bị bỏ lỡ sang danh sách thông báo thông thường
	service, ok := c.service.(*AppointmentService)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}
	service.notificationManager.MoveMissedToNotifications(authPayload.Username)

	// Thời gian tối đa chờ đợi: 30 giây
	timeout := 30 * time.Second

	// Chờ thông báo từ NotificationManager
	notification, hasNotification := service.notificationManager.WaitForNotification(authPayload.Username, timeout)
	if !hasNotification {
		// Trả về mã trạng thái 204 (No Content) nếu không có thông báo trong thời gian chờ
		ctx.Status(http.StatusNoContent)
		return
	}

	// Trả về thông báo cho client
	ctx.JSON(http.StatusOK, notification)
}

// GetNotifications trả về danh sách tất cả thông báo hiện tại và thông báo bị bỏ lỡ của người dùng
func (c *AppointmentController) getNotifications(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	service, ok := c.service.(*AppointmentService)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// Chuyển các thông báo bị bỏ lỡ sang danh sách thông báo thông thường
	service.notificationManager.MoveMissedToNotifications(authPayload.Username)

	// Lấy danh sách thông báo từ NotificationManager
	notifications := service.notificationManager.GetNotifications(authPayload.Username)

	ctx.JSON(http.StatusOK, notifications)
}

// GetMissedNotifications trả về danh sách thông báo bị bỏ lỡ của người dùng
func (c *AppointmentController) getMissedNotifications(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Lấy danh sách thông báo bị bỏ lỡ từ NotificationManager
	missed := c.service.(*AppointmentService).notificationManager.GetMissedNotifications(authPayload.Username)

	ctx.JSON(http.StatusOK, missed)
}

// RegisterUserRole đăng ký vai trò cho người dùng để nhận thông báo phù hợp
func (c *AppointmentController) registerUserRole(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := db.StoreDB.GetUser(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	// Kiểm tra tính hợp lệ của vai trò
	validRoles := map[string]bool{
		"admin":  true,
		"doctor": true,
		"nurse":  true,
	}

	role := user.Role.String

	if !validRoles[role] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Đăng ký vai trò với NotificationManager
	c.service.(*AppointmentService).notificationManager.SetUserRole(authPayload.Username, role)

	ctx.JSON(http.StatusOK, gin.H{"message": "Role registered successfully"})
}

// GetNotificationsFromDB lấy danh sách thông báo từ cơ sở dữ liệu
func (c *AppointmentController) getNotificationsFromDB(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	offset := (pagination.Page - 1) * pagination.PageSize
	limit := pagination.PageSize

	// Lấy thông báo từ database
	service, ok := c.service.(*AppointmentService)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	dbNotifications, err := service.notificationManager.GetNotificationsFromDB(ctx, authPayload.Username, int32(limit), int32(offset))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Không thể lấy thông báo từ database: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, dbNotifications)
}

// MarkNotificationAsRead đánh dấu thông báo đã đọc
func (c *AppointmentController) markNotificationAsRead(ctx *gin.Context) {
	// Xác thực người dùng
	_, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Lấy ID thông báo
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Đánh dấu thông báo đã đọc
	service, ok := c.service.(*AppointmentService)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	err = service.notificationManager.MarkNotificationAsReadInDB(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Không thể đánh dấu thông báo đã đọc: %v", err)})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Đã đánh dấu thông báo đã đọc"})
}

// MarkAllNotificationsAsRead đánh dấu tất cả thông báo của người dùng đã đọc
func (c *AppointmentController) markAllNotificationsAsRead(ctx *gin.Context) {
	// Xác thực người dùng
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Xóa tất cả thông báo trong bộ nhớ
	service, ok := c.service.(*AppointmentService)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
		return
	}

	// Xóa tất cả thông báo trong database
	err = service.notificationManager.DeleteNotificationsByUsernameFromDB(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Không thể xóa thông báo từ database: %v", err)})
		return
	}

	// Xóa thông báo trong bộ nhớ
	service.notificationManager.ClearNotifications(authPayload.Username)
	service.notificationManager.ClearMissedNotifications(authPayload.Username)

	ctx.JSON(http.StatusOK, gin.H{"message": "Đã đánh dấu tất cả thông báo đã đọc"})
}

func (c *AppointmentController) getCompletedAppointments(ctx *gin.Context) {
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	appointments, err := c.service.GetCompletedAppointments(ctx, authPayload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get completed appointments"})
		return
	}

	ctx.JSON(http.StatusOK, appointments)
}
