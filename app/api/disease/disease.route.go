package disease

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
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
			},
		},
	}

	{

		authRoute.GET("/", diceaseApi.controller.getDiceaseAnhMedicinesInfo)
		authRoute.GET("/treatment/:disease_id", diceaseApi.controller.getTreatmentByDiseaseId)
	}
	{
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment", diceaseApi.controller.CreateTreatment)
		perRoute([]perms.Permission{perms.ManageDisease, perms.ManageTreatment}).POST("/treatment/:treatment_id/phase", diceaseApi.controller.CreateTreatmentPhase)
	}

}
