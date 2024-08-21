package user

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	user := routerGroup.RouterDefault.Group("/user")
	authRoute := routerGroup.RouterAuth(user)

	// Khoi tao api
	userApi := &UserApi{
		&UserController{
			service: &UserService{
				storeDB : db.Store,
			},
		},
	}

	{
		authRoute.POST("/create", userApi.controller.createUser)
	}

}