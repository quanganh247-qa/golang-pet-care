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
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.PUT("/:appointment_id", appointmentApi.controller.updateAppointmentStatus)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsOfDoctor)
<<<<<<< HEAD
=======
>>>>>>> 323513c (appointment api)
=======
		authRoute.PUT("/:appointment_id", appointmentApi.controller.updateAppointmentStatus)
>>>>>>> 7cfffa9 (update dtb and appointment)
=======
>>>>>>> 4b8e9b6 (update appointment api)

	}

}
