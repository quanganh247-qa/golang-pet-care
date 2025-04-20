package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup, taskDistributor worker.TaskDistributor, ws *websocket.WSClientManager) {
	medicine := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(medicine)
	permsRoute := routerGroup.RouterPermission(medicine)
	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	MedicineApi := &MedicineApi{
		&MedicineController{
			service: &MedicineService{
				storeDB:         db.StoreDB,
				taskDistributor: taskDistributor,
				ws:              ws,
			},
		},
	}

	{
		authRoute.POST("/medicine", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/medicine/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/medicines/:pet_id", MedicineApi.controller.ListMedicines)
		authRoute.PUT("/medicine/:medicine_id", MedicineApi.controller.UpdateMedicine)
		authRoute.GET("/medicines", MedicineApi.controller.GetAllMedicines)
		authRoute.GET("/medicines/count", MedicineApi.controller.CountAllMedicines)

		authRoute.GET("/medicine/ws", MedicineApi.controller.HandleWebSocket)

		medicineInventoryGroup := permsRoute([]perms.Permission{perms.ManageTreatment})

		medicineInventoryGroup.POST("/medicine/transaction", MedicineApi.controller.CreateMedicineTransaction)
		medicineInventoryGroup.GET("/medicine/transactions", MedicineApi.controller.GetMedicineTransactions)
		medicineInventoryGroup.GET("/medicine/:medicine_id/transactions", MedicineApi.controller.GetMedicineTransactionsByMedicineID)

		medicineInventoryGroup.POST("/medicine/supplier", MedicineApi.controller.CreateSupplier)
		medicineInventoryGroup.GET("/medicine/supplier/:supplier_id", MedicineApi.controller.GetSupplierByID)
		medicineInventoryGroup.GET("/medicine/suppliers", MedicineApi.controller.GetAllSuppliers)
		medicineInventoryGroup.PUT("/medicine/supplier/:supplier_id", MedicineApi.controller.UpdateSupplier)

		medicineInventoryGroup.GET("/medicine/alerts/expiring", MedicineApi.controller.GetExpiringMedicines)
		medicineInventoryGroup.GET("/medicine/alerts/lowstock", MedicineApi.controller.GetLowStockMedicines)
	}
}
