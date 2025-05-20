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

	// Test/Vaccine API routes (combined functionality)
	{
		// Generic item routes (handle both tests and vaccines)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/items", controller.ListAllItems)
		permsRoute([]perms.Permission{perms.ManageTest}).POST("/test-orders", controller.CreateOrder)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/item/:id", controller.GetItemByID)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test-orders", controller.GetOrderedItemsByAppointment)
		permsRoute([]perms.Permission{perms.ManageTest}).DELETE("/item/:item_id", controller.SoftDeleteItem)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test-categories", controller.ListTestCategories)
		// Vaccine-specific routes
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/vaccines", controller.ListVaccines)
	}

	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/appointments/test-orders", controller.GetAllAppointmentsWithOrders)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/tests", controller.GetTestByAppointment)
	}
}
