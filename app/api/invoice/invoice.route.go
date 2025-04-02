package invoice

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	test := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(test)
	// permsRoute := routerGroup.RouterPermission(test)

	controller := NewInvoiceController(db.StoreDB)

	authRoute.POST("/invoice", controller.CreateInvoice)
	authRoute.GET("/invoice/:id", controller.GetInvoiceByID)
	authRoute.GET("/invoices", controller.ListInvoices)
	authRoute.PUT("/invoice/:id/status", controller.UpdateInvoiceStatus)
	authRoute.DELETE("/invoice/:id", controller.DeleteInvoice)
	authRoute.PUT("/invoice/:id/items/:item_id", controller.UpdateInvoiceItem)
	authRoute.DELETE("/invoice/:id/items/:item_id", controller.DeleteInvoiceItem)
	// Test order invoice route
	authRoute.POST("/invoices/test-order/:test_order_id", controller.CreateInvoiceFromTestOrder)

}
