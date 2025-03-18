<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("/", doctorApi.controller.getAllDoctor)
		authRoute.GET("/:doctor_id", doctorApi.controller.getDoctorById)
<<<<<<< HEAD
		authRoute.GET("/profile", doctorApi.controller.getDoctorProfile)
	}
	{
=======

=======
		authRoute.GET("/", doctorApi.controller.getAllDoctor)
		authRoute.GET("/:doctor_id", doctorApi.controller.getDoctorById)
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
		// Protected routes (require authentication)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
		authRoute.GET("/profile", doctorApi.controller.getDoctorProfile)
	}
	{
<<<<<<< HEAD
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/", doctorApi.controller.getAllDoctor)
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> 71b74e9 (feat(appointment): add room management and update appointment functionality.)
=======

		// Protected routes (require authentication)
		authRoute.GET("/profile", doctorApi.controller.getDoctorProfile)
	}
	{
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/", doctorApi.controller.getAllDoctor)
>>>>>>> 4ccd381 (Update appointment flow)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/shifts", doctorApi.controller.getShifts)
		perRoute([]perms.Permission{perms.ManageDoctor}).POST("/shifts", doctorApi.controller.createShift)
	}
}
<<<<<<< HEAD
<<<<<<< HEAD
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
package doctor
>>>>>>> ffc9071 (AI suggestion)
=======
package doctor
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)
