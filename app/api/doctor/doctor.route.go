package doctor

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	doctor := routerGroup.RouterDefault.Group("/doctor")
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
		doctor.POST("/login", doctorApi.controller.loginDoctor)

		// Protected routes (require authentication)
		authRoute.GET("/profile", doctorApi.controller.getDoctorProfile)
	}
	{
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/", doctorApi.controller.getAllDoctor)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/shifts", doctorApi.controller.getShifts)
		perRoute([]perms.Permission{perms.ManageDoctor}).POST("/shifts", doctorApi.controller.createShift)
	}
}
