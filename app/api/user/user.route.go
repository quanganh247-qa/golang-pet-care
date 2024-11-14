package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, config util.Config) {
	user := routerGroup.RouterDefault.Group("/user")
	authRoute := routerGroup.RouterAuth(user)

	// Khoi tao api
	userApi := &UserApi{
		&UserController{
			service: &UserService{
				storeDB:         db.StoreDB, // This should refer to the actual instance
				redis:           redis.Client,
				taskDistributor: taskDistributor,
				config:          config,
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
<<<<<<< HEAD
		authRoute.PUT("/", userApi.controller.updatetUser)
		authRoute.PUT("/avatar", userApi.controller.updatetUserAvatar)
		user.PUT("/reset-password", userApi.controller.ForgotPassword)
		authRoute.PUT("/change-password", userApi.controller.UpdatePassword)
=======
>>>>>>> bed48f6 (route logout)

<<<<<<< HEAD
		user.GET("/sessioninfo", userApi.controller.sessioninfo)
		user.GET("/userinfo", userApi.controller.userinfo)
		authRoute.GET("/roles", userApi.controller.GetAllRole)
=======
		// Doctor
		authRoute.POST("/create-doctor", userApi.controller.createDoctor)
		authRoute.POST("/timeslots", userApi.controller.insertTimeSlots)
		authRoute.GET("/doctor/:id", userApi.controller.getDoctor)

		// Schedule
		authRoute.GET("/time-slots/:doctor_id", userApi.controller.getTimeSlots)
		authRoute.GET("/all-time-slots/:doctor_id", userApi.controller.getAllTimeSlots)
		authRoute.PUT("/update-available/:id", userApi.controller.updateDoctorAvailableTime)

		// Token info for google calendar
		// authRoute.POST("/token-info", userApi.controller.insertTokenInfo)
>>>>>>> 79a3bcc (medicine api)

	}

}
