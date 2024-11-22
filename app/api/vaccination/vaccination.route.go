package vaccination

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {

	vaccination := routerGroup.RouterDefault.Group("/vaccination")
	authRoute := routerGroup.RouterAuth(vaccination)

	vaccinationController := &VaccinationController{
		service: &VaccinationService{
			storeDB: db.StoreDB,
		},
	}
	{
		authRoute.POST("/create", vaccinationController.CreateVaccination)
		authRoute.GET("/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("/:vaccination_id", vaccinationController.DeleteVaccination)
	}
}
