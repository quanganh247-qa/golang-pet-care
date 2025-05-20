package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

// Routes sets up all appointment related routes
// @Summary Appointment routes setup
// @Description Initializes all appointment related API endpoints
// @Tags appointments
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, ws *websocket.WSClientManager) {
	appointment := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(appointment)

	// Khởi tạo NotificationManager cho Long polling
	notificationManager := NewNotificationManager()

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
				storeDB:             db.StoreDB,
				taskDistributor:     taskDistributor,
				ws:                  ws,
				notificationManager: notificationManager,
			},
		},
	}

	{
		authRoute.POST("appointment/", appointmentApi.controller.createAppointment)

		authRoute.POST("appointment/confirm/:id", appointmentApi.controller.confirmAppointment)

		authRoute.POST("appointment/check-in/:id", appointmentApi.controller.checkinAppointment)

		authRoute.GET("appointment/user", appointmentApi.controller.getAppointmentsByUser)

		authRoute.GET("appointment/:id", appointmentApi.controller.getAppointmentByID)

		authRoute.DELETE("appointment/:id", appointmentApi.controller.deleteAppointment)

		authRoute.GET("appointment/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)

		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)

		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)

		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)

		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)

		authRoute.GET("appointments/history", appointmentApi.controller.getCompletedAppointments)

		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)

		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)

		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)

		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)

		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)

		authRoute.GET("appointments/statistic", appointmentApi.controller.GetAppointmentDistribution)

		// Walk-in appointments
		authRoute.POST("/appointments/walk-in", appointmentApi.controller.CreateWalkInAppointment)

		// WebSocket routes
		authRoute.GET("/appointment/websocket", appointmentApi.controller.HandleWebSocket)

		authRoute.GET("/appointment/state", appointmentApi.controller.GetAppointmentByState)

		// Long polling routes
		authRoute.GET("/appointment/notifications/wait", appointmentApi.controller.waitForNotifications)
		authRoute.GET("/appointment/notifications", appointmentApi.controller.getNotifications)

		// Database notification routes
		authRoute.GET("/appointment/notifications/db", appointmentApi.controller.getNotificationsFromDB)
		authRoute.PUT("/appointment/notifications/read/:id", appointmentApi.controller.markNotificationAsRead)
		authRoute.PUT("/appointment/notifications/read-all", appointmentApi.controller.markAllNotificationsAsRead)

		// Đăng ký vai trò người dùng
		authRoute.POST("/appointment/register-role", appointmentApi.controller.registerUserRole)
	}
}
