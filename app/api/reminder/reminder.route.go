package reminder

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {

	// login -> middleware (check token ) -> controller -> service -> model -> db
	// http://localhost:8080/api/v1/remider (jwt))
	//token

	reminder := routerGroup.RouterDefault.Group("/reminder")
	authRoute := routerGroup.RouterAuth(reminder)

	reminderAPI := &ReminderApi{
		&ReminderController{
			reminderService: &RemiderService{
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.POST("/create", reminderAPI.controller.CreateReminderController)
	}

	// api -> controller -> service -> model -> db

}
