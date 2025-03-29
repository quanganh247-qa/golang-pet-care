package doctor

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	doctor := routerGroup.RouterDefault.Group("/doctor")
	authRoute := routerGroup.RouterAuth(doctor)
	// perRoute := routerGroup.RouterPermission(doctor)

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
		doctor.POST("/login", doctorApi.controller.loginDoctor)
		authRoute.GET("/", doctorApi.controller.getAllDoctor)
		authRoute.GET("/:doctor_id", doctorApi.controller.getDoctorById)
		authRoute.GET("/profile", doctorApi.controller.getDoctorProfile)
		// Private routes
		authRoute.GET("/shifts", doctorApi.controller.getShifts)
		authRoute.POST("/shifts", doctorApi.controller.createShift)

	}
	// {
	// 	perRoute([]perms.Permission{perms.ManageDoctor}).GET("/shifts", doctorApi.controller.getShifts)
	// 	perRoute([]perms.Permission{perms.ManageDoctor}).POST("/shifts", doctorApi.controller.createShift)
	// }
}
