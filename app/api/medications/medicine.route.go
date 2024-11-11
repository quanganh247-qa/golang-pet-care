package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	medicine := routerGroup.RouterDefault.Group("/medicine")
	authRoute := routerGroup.RouterAuth(medicine)
	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	MedicineApi := &MedicineApi{
		&MedicineController{
			service: &MedicineService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		authRoute.POST("/create", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/list/:pet_id", MedicineApi.controller.ListMedicines)
		// authRoute.GET("/", MedicineApi.controller.ListMedicinesByUsername)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
		// authRoute.DELETE("/delete/:Medicineid", MedicineApi.controller.DeleteMedicine)
	}

}
