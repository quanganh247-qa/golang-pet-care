package medications

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService) {
	medicine := routerGroup.RouterDefault.Group("/medicine")
	authRoute := routerGroup.RouterAuth(medicine)
	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	MedicineApi := &MedicineApi{
		&MedicineController{
			service: &MedicineService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				es:      es,
			},
		},
	}

	{
		authRoute.POST("/", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/medicines/:pet_id", MedicineApi.controller.ListMedicines)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
	}

}
