package disease

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup, es *elasticsearch.ESService) {
<<<<<<< HEAD
	dicease := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(dicease)
	perRoute := routerGroup.RouterPermission(dicease)

	// Khoi tao api
	diseaseApi := &DiseaseApi{
		&DiseaseController{
			service: &DiseaseService{
				storeDB: db.StoreDB, // This should refer to the actual instance
				es:      es,
=======
=======
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
>>>>>>> 3bf345d (happy new year)
)

func Routes(routerGroup middleware.RouterGroup) {
=======
>>>>>>> e859654 (Elastic search)
	dicease := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(dicease)
	perRoute := routerGroup.RouterPermission(dicease)

	// Khoi tao api
	diseaseApi := &DiseaseApi{
		&DiseaseController{
			service: &DiseaseService{
				storeDB: db.StoreDB, // This should refer to the actual instance
<<<<<<< HEAD
>>>>>>> 6c35562 (dicease and treatment plan)
=======
				es:      es,
>>>>>>> e859654 (Elastic search)
=======
=======
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
>>>>>>> 3bf345d (happy new year)
)

func Routes(routerGroup middleware.RouterGroup) {
	dicease := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(dicease)
	perRoute := routerGroup.RouterPermission(dicease)

	// Khoi tao api
	diceaseApi := &DiceaseApi{
		&DiceaseController{
			service: &DiceaseService{
				storeDB: db.StoreDB, // This should refer to the actual instance
>>>>>>> 6c35562 (dicease and treatment plan)
			},
		},
	}

	{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase", diseaseApi.controller.CreateTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase/:phase_id/medicine", diseaseApi.controller.AssignMedicineToTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageTreatment}).PUT("/treatment/:treatment_id/phase/:phase_id", diseaseApi.controller.UpdateTreatmentPhaseStatus)
		perRoute([]perms.Permission{perms.ManageTreatment}).GET("/treatment/:treatment_id/prescription", diseaseApi.controller.GenerateMedicineOnlyPrescription)
	}

	dicease.POST("/disease", diseaseApi.controller.CreateDisease)
<<<<<<< HEAD
=======
		authRoute.GET("/", diceaseApi.controller.getDiceaseAnhMedicinesInfo)
		authRoute.GET("/treatment-plan", diceaseApi.controller.getDiseaseTreatmentPlanWithPhases)
=======
>>>>>>> 3bf345d (happy new year)

		// authRoute.GET("/", diceaseApi.controller.getDiceaseAnhMedicinesInfo)
		// authRoute.GET("/treatment/:disease_id", diceaseApi.controller.getTreatmentByDiseaseId)
		authRoute.GET("/pet/:pet_id/treatments", diceaseApi.controller.GetTreatmentsByPetID)
		authRoute.GET("/treatment/:treatment_id/phases", diceaseApi.controller.GetTreatmentPhasesByTreatmentID)
		authRoute.GET("/treatment/:treatment_id/phases/:phase_id/medicines", diceaseApi.controller.GetMedicinesByPhaseID)
		// for pet owner

=======
		authRoute.GET("/pet/:pet_id/treatments", diseaseApi.controller.GetTreatmentsByPetID)
		authRoute.GET("/treatment/:treatment_id/phases", diseaseApi.controller.GetTreatmentPhasesByTreatmentID)
		authRoute.GET("/treatment/:treatment_id/phases/:phase_id/medicines", diseaseApi.controller.GetMedicinesByPhaseID)
>>>>>>> e859654 (Elastic search)
	}
<<<<<<< HEAD
>>>>>>> 6c35562 (dicease and treatment plan)
=======
	{
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment", diseaseApi.controller.CreateTreatment)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase", diseaseApi.controller.CreateTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase/:phase_id/medicine", diseaseApi.controller.AssignMedicineToTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageTreatment}).PUT("/treatment/:treatment_id/phase/:phase_id", diseaseApi.controller.UpdateTreatmentPhaseStatus)
	}
>>>>>>> 3bf345d (happy new year)

	dicease.POST("/disease", diseaseApi.controller.CreateDisease)
=======

>>>>>>> ada3717 (Docker file)
=======
		authRoute.GET("/", diceaseApi.controller.getDiceaseAnhMedicinesInfo)
		authRoute.GET("/treatment-plan", diceaseApi.controller.getDiseaseTreatmentPlanWithPhases)
=======
>>>>>>> 3bf345d (happy new year)

		// authRoute.GET("/", diceaseApi.controller.getDiceaseAnhMedicinesInfo)
		// authRoute.GET("/treatment/:disease_id", diceaseApi.controller.getTreatmentByDiseaseId)
		authRoute.GET("/pet/:pet_id/treatments", diceaseApi.controller.GetTreatmentsByPetID)
		authRoute.GET("/treatment/:treatment_id/phases", diceaseApi.controller.GetTreatmentPhasesByTreatmentID)

	}
	{
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment", diceaseApi.controller.CreateTreatment)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase", diceaseApi.controller.CreateTreatmentPhase)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase/:phase_id/medicine", diceaseApi.controller.AssignMedicineToTreatmentPhase)
	}

>>>>>>> 6c35562 (dicease and treatment plan)
}
