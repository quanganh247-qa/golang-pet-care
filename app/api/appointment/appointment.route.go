package appointment

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	appointment := routerGroup.RouterDefault.Group("/")
=======
=======
>>>>>>> e859654 (Elastic search)
	appointment := routerGroup.RouterDefault.Group("/appointment")
>>>>>>> e859654 (Elastic search)
=======
	appointment := routerGroup.RouterDefault.Group("/")
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
	appointment := routerGroup.RouterDefault.Group("/")
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
	authRoute := routerGroup.RouterAuth(appointment)

	// Khoi tao api
	appointmentApi := &AppointmentApi{
		&AppointmentController{
			service: &AppointmentService{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
=======
				storeDB: db.StoreDB,
>>>>>>> 685da65 (latest update)
=======
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
>>>>>>> e859654 (Elastic search)
=======
				storeDB: db.StoreDB,
>>>>>>> 685da65 (latest update)
=======
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
>>>>>>> e859654 (Elastic search)
			},
		},
	}

	{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
		authRoute.POST("appointment/", appointmentApi.controller.createAppointment)
		authRoute.POST("appointment/confirm/:id", appointmentApi.controller.confirmAppointment)
		authRoute.POST("appointment/check-in/:id", appointmentApi.controller.checkinAppointment)
		authRoute.GET("appointment/user", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("appointment/:id", appointmentApi.controller.getAppointmentByID)
		authRoute.GET("appointment/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)
		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)
		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)
=======
		authRoute.POST("/create", appointmentApi.controller.createAppointment)
=======
		authRoute.POST("/", appointmentApi.controller.createAppointment)
<<<<<<< HEAD
<<<<<<< HEAD
>>>>>>> 685da65 (latest update)
		authRoute.PUT("/:appointment_id", appointmentApi.controller.updateAppointmentStatus)
<<<<<<< HEAD
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("/:appointment_id", appointmentApi.controller.getAppointmentByID)
>>>>>>> 7e35c2e (get appointment detail)
=======
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsOfDoctor)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByPetOfUser)
=======
		authRoute.POST("/", appointmentApi.controller.createAppointment)
		authRoute.POST("/confirm/:id", appointmentApi.controller.confirmAppointment)
		authRoute.POST("/check-in/:id", appointmentApi.controller.checkinAppointment)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
<<<<<<< HEAD
>>>>>>> 685da65 (latest update)
		authRoute.GET("/:appointment_id", appointmentApi.controller.getAppointmentByID)
<<<<<<< HEAD
>>>>>>> 7e35c2e (get appointment detail)

		// soap
		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)
		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)
		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)
=======
		authRoute.POST("/confirm/:appointment_id", appointmentApi.controller.confirmAppointment)
=======
		authRoute.POST("/confirm/:id", appointmentApi.controller.confirmAppointment)
<<<<<<< HEAD
>>>>>>> e859654 (Elastic search)
=======
		authRoute.POST("/check-in/:id", appointmentApi.controller.checkinAppointment)
>>>>>>> 4ccd381 (Update appointment flow)
		authRoute.GET("/", appointmentApi.controller.getAppointmentsByUser)
		authRoute.GET("/:id", appointmentApi.controller.getAppointmentByID)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)
		authRoute.GET("/all", appointmentApi.controller.getAllAppointments)

<<<<<<< HEAD
=======
		// authRoute.GET("/all", appointmentApi.controller.getAllAppointments)
		authRoute.GET("appointments/", appointmentApi.controller.getAllAppointmentsByDate)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======
>>>>>>> dc47646 (Optimize SQL query)
=======
		// authRoute.GET("/all", appointmentApi.controller.getAllAppointments)
		authRoute.GET("appointments/", appointmentApi.controller.getAllAppointmentsByDate)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
		authRoute.GET("/doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======

		// soap
		authRoute.POST("/soap/:appointment_id", appointmentApi.controller.createSOAP)
		authRoute.PUT("/soap/:appointment_id", appointmentApi.controller.updateSOAP)
>>>>>>> e859654 (Elastic search)
=======
		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)
		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)
		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)

		// soap
		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)
		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)
		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
=======
		authRoute.GET("/:id", appointmentApi.controller.getAppointmentByID)
>>>>>>> e859654 (Elastic search)
		authRoute.GET("/doctor/:doctor_id", appointmentApi.controller.getAppointmentsByDoctor)
		authRoute.GET("/all", appointmentApi.controller.getAllAppointments)

		// time slot
		authRoute.GET("/doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
<<<<<<< HEAD
>>>>>>> b393bb9 (add service and add permission)
=======

		// soap
		authRoute.POST("/soap/:appointment_id", appointmentApi.controller.createSOAP)
		authRoute.PUT("/soap/:appointment_id", appointmentApi.controller.updateSOAP)
>>>>>>> e859654 (Elastic search)
=======
		authRoute.GET("appointments", appointmentApi.controller.getAllAppointments)
		authRoute.PUT("appointment/:id", appointmentApi.controller.updateAppointment)
		authRoute.GET("doctor/:doctor_id/time-slot", appointmentApi.controller.getAvailableTimeSlots)
		authRoute.GET("appointments/pet/:pet_id/history", appointmentApi.controller.getHistoryAppointmentsByPetID)
		authRoute.GET("appointments/queue", appointmentApi.controller.getQueue)
		authRoute.PUT("appointments/queue/:id/status", appointmentApi.controller.updateQueueItemStatus)

		// soap
		authRoute.POST("appointment/:id/soap", appointmentApi.controller.createSOAP)
		authRoute.PUT("appointment/:id/soap", appointmentApi.controller.updateSOAP)
		authRoute.GET("appointment/:id/soap", appointmentApi.controller.getSOAPByAppointmentID)
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
	}

}
