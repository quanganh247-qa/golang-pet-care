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
		authRoute.GET("/doctors", doctorApi.controller.getAllDoctor)
		authRoute.GET("/doctor/:doctor_id", doctorApi.controller.getDoctorById)
		authRoute.GET("/doctor/profile", doctorApi.controller.getDoctorProfile)
		authRoute.PUT("/doctor/profile", doctorApi.controller.editDoctorProfile)
		// Private routes
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/shifts", doctorApi.controller.getShifts)
		perRoute([]perms.Permission{perms.ManageDoctor}).POST("/doctor/shifts", doctorApi.controller.createShift)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/:doctor_id/shifts", doctorApi.controller.getShiftByDoctorId)
		perRoute([]perms.Permission{perms.ManageDoctor}).POST("/doctor/:doctor_id/leave", doctorApi.controller.createLeaveRequest)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/:doctor_id/leave", doctorApi.controller.getLeaveRequests)
		perRoute([]perms.Permission{perms.ManageDoctor}).PUT("/leave/:leave_id", doctorApi.controller.updateLeaveRequest)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/:doctor_id/attendance", doctorApi.controller.getDoctorAttendance)
		perRoute([]perms.Permission{perms.ManageDoctor}).GET("/doctor/:doctor_id/workload", doctorApi.controller.getDoctorWorkload)
	}
}
