package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
)

func Routes(routerGroup middleware.RouterGroup) {
	user := routerGroup.RouterDefault.Group("/user")
	authRoute := routerGroup.RouterAuth(user)
	// user.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	userApi := &UserApi{
		&UserController{
			service: &UserService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				redis:   redis.Client,
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		user.POST("/create", userApi.controller.createUser)
		user.GET("/all", userApi.controller.getAllUsers)
		authRoute.GET("/", userApi.controller.getUserDetails)
		user.POST("/login", userApi.controller.loginUser)
		user.PUT("/verify-email", userApi.controller.verifyEmail)
		authRoute.GET("/refresh_token", userApi.controller.getAccessToken)
		authRoute.POST("/logout", userApi.controller.logoutUser)

		// Doctor
		authRoute.POST("/create-doctor", userApi.controller.createDoctor)
		authRoute.POST("/timeslots", userApi.controller.insertTimeSlots)
		authRoute.GET("/doctor/:id", userApi.controller.getDoctor)

		// Schedule
		authRoute.GET("/time-slots/:doctor_id", userApi.controller.getTimeSlots)
		authRoute.GET("/all-time-slots/:doctor_id", userApi.controller.getAllTimeSlots)
		authRoute.PUT("/update-available/:id", userApi.controller.updateDoctorAvailableTime)

	}

}
