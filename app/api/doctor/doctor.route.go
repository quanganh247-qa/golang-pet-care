package doctor

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	doctor := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(doctor)
	perRoute := routerGroup.RouterPermission(doctor)

	// Initialize API
	doctorApi := &DoctorApi{
		controller: &DoctorController{
			service: &DoctorService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		// Public routes
		doctor.POST("/doctor/login", doctorApi.controller.loginDoctor)
		authRoute.GET("/doctor", doctorApi.controller.getAllDoctor)
		authRoute.GET("/doctor/:doctor_id", doctorApi.controller.getDoctorById)
		authRoute.GET("/doctor/profile", doctorApi.controller.getDoctorProfile)
		// Private routes
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/shifts", doctorApi.controller.getShifts)
		perRoute([]perms.Permission{perms.ManageDoctor}).POST("/doctor/shifts", doctorApi.controller.createShift)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/:doctor_id/shifts", doctorApi.controller.getShiftByDoctorId)

	}
}
