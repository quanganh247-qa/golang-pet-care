package test

import (
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService, ws *websocket.WSClientManager) {
	test := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(test)
	permsRoute := routerGroup.RouterPermission(test)
	controller := NewTestController(es, ws)

	// authRoute.POST("/", controller.CreateTest)
	// authRoute.PUT("/:test_id/status", controller.UpdateTestStatus)
	// authRoute.POST("/:test_id/results", controller.AddTestResult)
	// authRoute.GET("/", controller.GetTestsByPetsID)
	// authRoute.GET("/:test_id/history", controller.GetStatusHistory)
	authRoute.GET("/ws", controller.HandleWebSocket)

	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/tests", controller.ListTests)
		permsRoute([]perms.Permission{perms.ManageTest}).POST("/test-orders", controller.CreateTestOrder)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test/:id", controller.GetTestByID)
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/test-orders", controller.GetOrderedTestsByAppointment)

		// permsRoute([]perms.Permission{perms.ManageTest}).GET("/test/:test_id", controller.UpdateTest)

	}
	{
		permsRoute([]perms.Permission{perms.ManageTest}).GET("/appointments/test-orders", controller.GetAllAppointmentsWithOrders)

	}
}
