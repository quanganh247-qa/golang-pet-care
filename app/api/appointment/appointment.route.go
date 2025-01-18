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
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.POST("/", appointmentApi.controller.createAppointment)
		authRoute.PUT("/:appointment_id", appointmentApi.controller.updateAppointmentStatus)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("/:appointment_id", appointmentApi.controller.getAppointmentByID)

	}

}
