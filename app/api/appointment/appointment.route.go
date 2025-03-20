package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

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
		// authRoute.GET("/all", appointmentApi.controller.getAllAppointments)
		authRoute.GET("appointments/", appointmentApi.controller.getAllAppointmentsByDate)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
		authRoute.GET("/doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)

		// soap
		authRoute.POST("/soap/:appointment_id", appointmentApi.controller.createSOAP)
		authRoute.PUT("/soap/:appointment_id", appointmentApi.controller.updateSOAP)
	}

}
