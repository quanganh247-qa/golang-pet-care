package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
=======
>>>>>>> 272832d (redis cache)
=======
>>>>>>> 272832d (redis cache)
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, config util.Config) {
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

<<<<<<< HEAD
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
>>>>>>> 6610455 (feat: redis queue)
=======
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, config util.Config) {
>>>>>>> 1a9e82a (reset password api)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

<<<<<<< HEAD
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor) {
>>>>>>> 6610455 (feat: redis queue)
=======
func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, config util.Config) {
>>>>>>> 1a9e82a (reset password api)
	user := routerGroup.RouterDefault.Group("/user")
	authRoute := routerGroup.RouterAuth(user)

	// Khoi tao api
	userApi := &UserApi{
		&UserController{
			service: &UserService{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
				storeDB:         db.StoreDB, // This should refer to the actual instance
				redis:           redis.Client,
				taskDistributor: taskDistributor,
				config:          config,
<<<<<<< HEAD
<<<<<<< HEAD
=======
				storeDB: db.StoreDB, // This should refer to the actual instance
				redis:   redis.Client,
<<<<<<< HEAD
>>>>>>> 272832d (redis cache)
=======
				storeDB:         db.StoreDB, // This should refer to the actual instance
				redis:           redis.Client,
				taskDistributor: taskDistributor,
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> 1a9e82a (reset password api)
=======
>>>>>>> 272832d (redis cache)
=======
				storeDB:         db.StoreDB, // This should refer to the actual instance
				redis:           redis.Client,
				taskDistributor: taskDistributor,
>>>>>>> 6610455 (feat: redis queue)
=======
>>>>>>> 1a9e82a (reset password api)
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		user.POST("/create", userApi.controller.createUser)
		user.GET("/all", userApi.controller.getAllUsers)
		authRoute.GET("/", userApi.controller.getUserDetails)
		user.POST("/login", userApi.controller.loginUser)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		user.POST("/verify_email", userApi.controller.verifyEmail)
		user.POST("/resend_otp/:username", userApi.controller.resendOTP)
=======
		user.PUT("/verify_email", userApi.controller.verifyEmail)
>>>>>>> 6610455 (feat: redis queue)
=======
		user.POST("/verify_email", userApi.controller.verifyEmail)
		user.POST("/resend_otp/:username", userApi.controller.resendOTP)
>>>>>>> edfe5ad (OTP verifycation)
=======
		user.PUT("/verify_email", userApi.controller.verifyEmail)
>>>>>>> 6610455 (feat: redis queue)
=======
		user.POST("/verify_email", userApi.controller.verifyEmail)
		user.POST("/resend_otp/:username", userApi.controller.resendOTP)
>>>>>>> edfe5ad (OTP verifycation)
		authRoute.GET("/refresh_token", userApi.controller.getAccessToken)
		authRoute.POST("/logout", userApi.controller.logoutUser)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.PUT("/", userApi.controller.updatetUser)
		authRoute.PUT("/avatar", userApi.controller.updatetUserAvatar)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		user.PUT("/reset-password", userApi.controller.ForgotPassword)
		authRoute.PUT("/change-password", userApi.controller.UpdatePassword)
=======
>>>>>>> bed48f6 (route logout)
=======
		authRoute.PUT("/", userApi.controller.updatetUser)
		authRoute.PUT("/avatar", userApi.controller.updatetUserAvatar)
>>>>>>> 473cd1d (uplaod image method)
=======
		user.PUT("/password", userApi.controller.ForgotPassword)
>>>>>>> 1a9e82a (reset password api)
=======
		user.PUT("/reset-password", userApi.controller.ForgotPassword)
		authRoute.PUT("/change-password", userApi.controller.UpdatePassword)
>>>>>>> a2c21c8 (update pass)
=======
>>>>>>> bed48f6 (route logout)
=======
		authRoute.PUT("/", userApi.controller.updatetUser)
		authRoute.PUT("/avatar", userApi.controller.updatetUserAvatar)
>>>>>>> 473cd1d (uplaod image method)
=======
		user.PUT("/password", userApi.controller.ForgotPassword)
>>>>>>> 1a9e82a (reset password api)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		user.GET("/sessioninfo", userApi.controller.sessioninfo)
		user.GET("/userinfo", userApi.controller.userinfo)
		authRoute.GET("/roles", userApi.controller.GetAllRole)
=======
		// Doctor
		authRoute.POST("/create-doctor", userApi.controller.createDoctor)
		authRoute.POST("/timeslots", userApi.controller.insertTimeSlots)
		authRoute.GET("/doctor/:id", userApi.controller.getDoctor)
		authRoute.GET("/doctors", userApi.controller.GetDoctors)

		// Schedule
		authRoute.GET("/time-slots/:doctor_id", userApi.controller.getTimeSlots)
		authRoute.GET("/all-time-slots/:doctor_id", userApi.controller.getAllTimeSlots)
		authRoute.PUT("/update-available/:id", userApi.controller.updateDoctorAvailableTime)

<<<<<<< HEAD
<<<<<<< HEAD
		// Token info for google calendar
		// authRoute.POST("/token-info", userApi.controller.insertTokenInfo)
<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 79a3bcc (medicine api)

=======
>>>>>>> e01abc5 (pet schedule api)
=======
>>>>>>> ae87825 (updated)
=======
		user.GET("/sessioninfo", userApi.controller.sessioninfo)
		user.GET("/userinfo", userApi.controller.userinfo)
		authRoute.GET("/roles", userApi.controller.GetAllRole)

>>>>>>> ada3717 (Docker file)
=======
>>>>>>> e01abc5 (pet schedule api)
	}

}
