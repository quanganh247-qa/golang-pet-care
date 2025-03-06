package medications

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService) {
=======
)

func Routes(routerGroup middleware.RouterGroup) {
>>>>>>> 79a3bcc (medicine api)
=======
=======
>>>>>>> a415f25 (new data)
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)
<<<<<<< HEAD

func Routes(routerGroup middleware.RouterGroup) {
>>>>>>> 79a3bcc (medicine api)
	medicine := routerGroup.RouterDefault.Group("/medicine")
	authRoute := routerGroup.RouterAuth(medicine)
	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	MedicineApi := &MedicineApi{
		&MedicineController{
			service: &MedicineService{
				storeDB: db.StoreDB, // This should refer to the actual instance
<<<<<<< HEAD
<<<<<<< HEAD
				es:      es,
=======
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 79a3bcc (medicine api)
			},
		},
	}

	{
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.POST("/", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/medicines/:pet_id", MedicineApi.controller.ListMedicines)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
=======
=======
>>>>>>> 79a3bcc (medicine api)
		authRoute.POST("/create", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/list/:pet_id", MedicineApi.controller.ListMedicines)
		// authRoute.GET("/", MedicineApi.controller.ListMedicinesByUsername)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
		// authRoute.DELETE("/delete/:Medicineid", MedicineApi.controller.DeleteMedicine)
<<<<<<< HEAD
>>>>>>> 79a3bcc (medicine api)
	}

}
=======
// func Routes(routerGroup middleware.RouterGroup) {
// 	medicine := routerGroup.RouterDefault.Group("/medicine")
// 	authRoute := routerGroup.RouterAuth(medicine)
// 	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())
=======
import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
)
>>>>>>> a415f25 (new data)

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

<<<<<<< HEAD
// }
>>>>>>> 6c35562 (dicease and treatment plan)
=======
	{
		authRoute.POST("/", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/medicines/:pet_id", MedicineApi.controller.ListMedicines)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
	}

}
>>>>>>> a415f25 (new data)
=======
	}

}
>>>>>>> 79a3bcc (medicine api)
=======
// func Routes(routerGroup middleware.RouterGroup) {
// 	medicine := routerGroup.RouterDefault.Group("/medicine")
// 	authRoute := routerGroup.RouterAuth(medicine)
// 	// Medicine.Use(middleware.IPbasedRateLimitingMiddleware())
=======
>>>>>>> a415f25 (new data)

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

<<<<<<< HEAD
// }
>>>>>>> 6c35562 (dicease and treatment plan)
=======
	{
		authRoute.POST("/", MedicineApi.controller.CreateMedicine)
		authRoute.GET("/:medicine_id", MedicineApi.controller.GetMedicineByID)
		authRoute.GET("/medicines/:pet_id", MedicineApi.controller.ListMedicines)
		// authRoute.GET("/", MedicineApi.controller.ListMedicinesByUsername)
		authRoute.PUT("/:medicine_id", MedicineApi.controller.UpdateMedicine)
		// authRoute.DELETE("/delete/:Medicineid", MedicineApi.controller.DeleteMedicine)
	}

}
>>>>>>> a415f25 (new data)
