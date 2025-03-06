package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
	appointment := routerGroup.RouterDefault.Group("/appointment")
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
		authRoute.POST("/", appointmentApi.controller.createAppointment)
		authRoute.POST("/confirm/:id", appointmentApi.controller.confirmAppointment)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("/:id", appointmentApi.controller.getAppointmentByID)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)
		authRoute.GET("/all", appointmentApi.controller.getAllAppointments)

		// time slot
		authRoute.GET("/doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)

		// soap
		authRoute.POST("/soap/:appointment_id", appointmentApi.controller.createSOAP)
		authRoute.PUT("/soap/:appointment_id", appointmentApi.controller.updateSOAP)
	}

}
