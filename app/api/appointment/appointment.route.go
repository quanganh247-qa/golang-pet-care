package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	appointment := routerGroup.RouterDefault.Group("/appointment")
	authRoute := routerGroup.RouterAuth(appointment)

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.POST("/", appointmentApi.controller.createAppointment)
		authRoute.POST("/confirm/:appointment_id", appointmentApi.controller.confirmAppointment)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("/:appointment_id", appointmentApi.controller.getAppointmentByID)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)
		// time slot
		authRoute.GET("/doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
	}

}
