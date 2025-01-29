package medical_records

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	MedicalRecord := routerGroup.RouterDefault.Group("/")
	// authRoute := routerGroup.RouterAuth(MedicalRecord)
	perRoute := routerGroup.RouterPermission(MedicalRecord)

	// Khoi tao api
	MedicalRecordApi := &MedicalRecordApi{
		&MedicalRecordController{
			service: &MedicalRecordService{
				storeDB: db.StoreDB, // This should refer to the actual instance
			},
		},
	}

	{
		// Example: Only users with the "MANAGE_SERVICES" permission can access this route
		perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-records", MedicalRecordApi.controller.CreateMedicalRecord)
		perRoute([]perms.Permission{perms.ManageMedicalRecords}).POST("/pets/:pet_id/medical-histories", MedicalRecordApi.controller.CreateMedicalHistory)
	}

}
