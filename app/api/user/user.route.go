package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
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
				// emailQueue: rabbitmq.Client.Email,
			},
		},
	}

	{
		// authRoute.POST("/create", userApi.controller.createUser)
		user.POST("/create", userApi.controller.createUser)
		user.GET("/all", userApi.controller.getAllUsers)
		user.POST("/login", userApi.controller.loginUser)
		user.PUT("/verify-email", userApi.controller.verifyEmail)
		authRoute.GET("/refresh_token", userApi.controller.getAccessToken)

		// Doctor
		authRoute.POST("/create-doctor", userApi.controller.createDoctor)
		authRoute.POST("/add-schedule", userApi.controller.addSchedule)
		authRoute.GET("/doctor/:id", userApi.controller.getDoctor)

	}

}
