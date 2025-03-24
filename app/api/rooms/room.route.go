package rooms

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	room := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(room)

	// Khoi tao api
	roomAPI := &RoomApi{
		&RoomController{
			service: &RoomService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	authRoute.POST("room/", roomAPI.controller.CreateRoom)
	// User routes
	authRoute.GET("rooms/", roomAPI.controller.ListRooms)
	authRoute.GET("room/:id", roomAPI.controller.GetRoomByID)

}
