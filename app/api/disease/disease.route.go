package disease

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService) {
	dicease := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(dicease)
	perRoute := routerGroup.RouterPermission(dicease)

	// Khoi tao api
	diseaseApi := &DiseaseApi{
		&DiseaseController{
			service: &DiseaseService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				es:      es,
			},
		},
	}

	{
		authRoute.GET("/pet/:pet_id/treatments", diseaseApi.controller.GetTreatmentsByPetID)
		authRoute.GET("/treatment/:treatment_id/phases", diseaseApi.controller.GetTreatmentPhasesByTreatmentID)
		authRoute.GET("/treatment/:treatment_id/phases/:phase_id/medicines", diseaseApi.controller.GetMedicinesByPhaseID)
	}
	{
		authRoute.GET("/pet/:pet_id/allergies", diseaseApi.controller.GetAllergiesByPetID)
		authRoute.POST("/pet/:pet_id/allerg", diseaseApi.controller.CreateAllergy)
	}
	{
		// perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).GET("/medicines", diseaseApi.controller.GetAllMedicines)
	}
	{
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment", diseaseApi.controller.CreateTreatment)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase", diseaseApi.controller.CreateTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase/:phase_id/medicine", diseaseApi.controller.AssignMedicineToTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageTreatment}).PUT("/treatment/:treatment_id/phase/:phase_id", diseaseApi.controller.UpdateTreatmentPhaseStatus)
		perRoute([]perms.Permission{perms.ManageTreatment}).GET("/treatment/:treatment_id/prescription", diseaseApi.controller.GenerateMedicineOnlyPrescription)
	}

	dicease.POST("/disease", diseaseApi.controller.CreateDisease)

}
