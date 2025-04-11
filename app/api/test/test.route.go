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

		// Vaccine-specific routes
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/vaccines", controller.ListVaccines)
	}

	// Legacy routes for backward compatibility
	// {
	// 	permsRoute([]perms.Permission{perms.ManageTest}).GET("/tests", controller.ListTests)
	// 	// permsRoute([]perms.Permission{perms.ManageTest}).POST("/test-orders", controller.CreateTestOrder)
	// 	permsRoute([]perms.Permission{perms.ManageTest}).GET("/test/:id", controller.GetTestByID)
	// 	// permsRoute([]perms.Permission{perms.ManageTest}).GET("/test-orders", controller.GetOrderedTestsByAppointment)
	// 	permsRoute([]perms.Permission{perms.ManageTest}).DELETE("/test/:test_id", controller.SoftDeleteTest)
	// }

	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/appointments/test-orders", controller.GetAllAppointmentsWithOrders)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/tests", controller.GetTestByAppointment)
	}
}
