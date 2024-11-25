package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
	user := routerGroup.RouterDefault.Group("/user")
	authRoute := routerGroup.RouterAuth(user)

	// Khoi tao api
	userApi := &UserApi{
		&UserController{
			service: &UserService{
				storeDB:         db.StoreDB, // This should refer to the actual instance
				redis:           redis.Client,
				taskDistributor: taskDistributor,
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		user.POST("/create", userApi.controller.createUser)
		user.GET("/all", userApi.controller.getAllUsers)
		authRoute.GET("/", userApi.controller.getUserDetails)
		user.POST("/login", userApi.controller.loginUser)
		user.POST("/verify_email", userApi.controller.verifyEmail)
		user.POST("/resend_otp/:username", userApi.controller.resendOTP)
		authRoute.GET("/refresh_token", userApi.controller.getAccessToken)
		authRoute.POST("/logout", userApi.controller.logoutUser)
		authRoute.PUT("/", userApi.controller.updatetUser)
		authRoute.PUT("/avatar", userApi.controller.updatetUserAvatar)

		// Doctor
		authRoute.POST("/create-doctor", userApi.controller.createDoctor)
		authRoute.POST("/timeslots", userApi.controller.insertTimeSlots)
		authRoute.GET("/doctor/:id", userApi.controller.getDoctor)
		authRoute.GET("/doctors", userApi.controller.GetDoctors)

		// Schedule
		authRoute.GET("/time-slots/:doctor_id", userApi.controller.getTimeSlots)
		authRoute.GET("/all-time-slots/:doctor_id", userApi.controller.getAllTimeSlots)
		authRoute.PUT("/update-available/:id", userApi.controller.updateDoctorAvailableTime)

	}

}
