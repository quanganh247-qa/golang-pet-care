package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	appointment := routerGroup.RouterDefault.Group("/appointment")
	authRoute := routerGroup.RouterAuth(appointment)
	// appointment.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				// emailQueue: rabbitmq.Client.Email,
			},
		},
	}

	{
		authRoute.POST("/create", appointmentApi.controller.createAppointment)
		authRoute.PUT("/:appointment_id", appointmentApi.controller.updateAppointmentStatus)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsOfDoctor)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByPetOfUser)
		authRoute.GET("/:appointment_id", appointmentApi.controller.getAppointmentByID)

	}

}
