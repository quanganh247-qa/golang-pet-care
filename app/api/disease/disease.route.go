package disease

import (
	"fmt"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/llm"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	dicease := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(dicease)
	perRoute := routerGroup.RouterPermission(dicease)

	// Initialize AIService with API key from config
	config := util.Configs
	aiService, err := llm.NewAIService(config.OpenAIAPIKey)
	if err != nil {
		// Log the error but continue with nil service
		// The AnalyzeSymptoms function will return an appropriate error
		fmt.Println("Error initializing AIService:", err)
	}

	// Khoi tao api
	diseaseApi := &DiseaseApi{
		&DiseaseController{
			service: &DiseaseService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				// es:      es,
			},
			llmService: *aiService, // Dereference the pointer
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
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment", diseaseApi.controller.CreateTreatment)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phases", diseaseApi.controller.CreateTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase/:phase_id/medicines", diseaseApi.controller.AssignMedicineToTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageTreatment}).PUT("/treatment/:treatment_id/phase/:phase_id", diseaseApi.controller.UpdateTreatmentPhaseStatus)
		perRoute([]perms.Permission{perms.ManageTreatment}).GET("/treatment/:treatment_id/prescription", diseaseApi.controller.GenerateMedicineOnlyPrescription)
	}

	dicease.POST("/disease", diseaseApi.controller.CreateDisease)

	dicease.POST("/symptom-analysis", diseaseApi.controller.AnalyzeSymptoms)

}
