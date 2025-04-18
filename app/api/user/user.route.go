package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, config util.Config) {
	user := routerGroup.RouterDefault.Group("/")
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
		user.POST("/user/create", userApi.controller.createUser)
		user.GET("/user/all", userApi.controller.getAllUsers)
		authRoute.GET("/user", userApi.controller.getUserDetails)
		user.POST("/user/login", userApi.controller.loginUser)
		user.POST("/user/verify_email", userApi.controller.verifyEmail)
		user.POST("/user/resend_otp/:username", userApi.controller.resendOTP)
		authRoute.GET("/user/refresh_token", userApi.controller.getAccessToken)
		authRoute.POST("/user/logout", userApi.controller.logoutUser)
		authRoute.PUT("/user", userApi.controller.updatetUser)
		authRoute.PUT("/user/avatar", userApi.controller.updatetUserAvatar)
		user.PUT("/user/reset-password", userApi.controller.ForgotPassword)
		authRoute.PUT("/user/change-password", userApi.controller.UpdatePassword)

		user.GET("/user/sessioninfo", userApi.controller.sessioninfo)
		user.GET("/user/userinfo", userApi.controller.userinfo)
		authRoute.GET("/user/roles", userApi.controller.GetAllRole)

		// New route for creating staff accounts - requires authentication
		authRoute.POST("/staff", userApi.controller.createStaff)
	}

}
