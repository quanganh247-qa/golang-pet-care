package test

import (
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup, ws *websocket.WSClientManager) {
	test := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(test)
	permsRoute := routerGroup.RouterPermission(test)
	controller := NewTestController(ws)

	authRoute.GET("/ws", controller.HandleWebSocket)

	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/tests", controller.ListTests)
		permsRoute([]perms.Permission{perms.ManageTest}).POST("/test-orders", controller.CreateTestOrder)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test/:id", controller.GetTestByID)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test-orders", controller.GetOrderedTestsByAppointment)
	}
	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/appointments/test-orders", controller.GetAllAppointmentsWithOrders)

	}
}
