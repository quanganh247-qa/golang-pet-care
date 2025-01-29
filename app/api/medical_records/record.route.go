package medical_records

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
=======
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
>>>>>>> 3bf345d (happy new year)
)

func Routes(routerGroup middleware.RouterGroup) {
	MedicalRecord := routerGroup.RouterDefault.Group("/")
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	authRoute := routerGroup.RouterAuth(MedicalRecord)
	// perRoute := routerGroup.RouterPermission(MedicalRecord)
=======
	// authRoute := routerGroup.RouterAuth(MedicalRecord)
	perRoute := routerGroup.RouterPermission(MedicalRecord)
>>>>>>> 3bf345d (happy new year)
=======
	authRoute := routerGroup.RouterAuth(MedicalRecord)
	// perRoute := routerGroup.RouterPermission(MedicalRecord)
>>>>>>> e859654 (Elastic search)
=======
	// authRoute := routerGroup.RouterAuth(MedicalRecord)
	perRoute := routerGroup.RouterPermission(MedicalRecord)
>>>>>>> 3bf345d (happy new year)

	// Khoi tao api
	MedicalRecordApi := &MedicalRecordApi{
		&MedicalRecordController{
			service: &MedicalRecordService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
		authRoute.GET("/pets/:pet_id/medical-records", MedicalRecordApi.controller.GetMedicalRecord)
		authRoute.POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
		authRoute.GET("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.ListMedicalHistory)
		authRoute.POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("/pets/:pet_id/medical-histories/:id", MedicalRecordApi.controller.GetMedicalHistoryByID)
	}

	// {
	// 	// Example: Only users with the "MANAGE_SERVICES" permission can access this route
	// 	perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
	// 	perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
	// }

=======
=======
>>>>>>> 3bf345d (happy new year)
		// Example: Only users with the "MANAGE_SERVICES" permission can access this route
		perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
		perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
	}

<<<<<<< HEAD
>>>>>>> 3bf345d (happy new year)
=======
		authRoute.POST("/pets/:pet_id/allergies", MedicalRecordApi.controller.CreateAllergy)
=======
>>>>>>> 4ccd381 (Update appointment flow)
		authRoute.GET("/pets/:pet_id/medical-histories/:id", MedicalRecordApi.controller.GetMedicalHistoryByID)
	}

	// {
	// 	// Example: Only users with the "MANAGE_SERVICES" permission can access this route
	// 	perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
	// 	perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
	// }

>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 3bf345d (happy new year)
}
