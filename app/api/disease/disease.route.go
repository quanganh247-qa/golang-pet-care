package disease

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	dicease := routerGroup.RouterDefault.Group("/disease")
	authRoute := routerGroup.RouterAuth(dicease)

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
		authRoute.GET("/treatment-plan", diceaseApi.controller.getDiseaseTreatmentPlanWithPhases)

	}

}
