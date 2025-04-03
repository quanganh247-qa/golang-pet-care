package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

// Routes sets up all appointment related routes
// @Summary Appointment routes setup
// @Description Initializes all appointment related API endpoints
// @Tags appointments
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
	appointment := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(appointment)

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
			},
		},
	}

	{

		authRoute.POST("appointment/", appointmentApi.controller.createAppointment)

		authRoute.POST("appointment/confirm/:id", appointmentApi.controller.confirmAppointment)

		authRoute.POST("appointment/check-in/:id", appointmentApi.controller.checkinAppointment)

		authRoute.GET("appointment/user", appointmentApi.controller.getAppointmentsByUser)

		authRoute.GET("appointment/:id", appointmentApi.controller.getAppointmentByID)

		authRoute.GET("appointment/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)

		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)

		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)

		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)

		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)

		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)

		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)

		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)

		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)

		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)

		authRoute.GET("appointments/statistic", appointmentApi.controller.GetAppointmentDistribution)
	}
}
