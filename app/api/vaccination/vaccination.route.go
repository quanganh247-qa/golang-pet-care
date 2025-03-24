package vaccination

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {

	vaccination := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(vaccination)

	vaccinationController := &VaccinationController{
		service: &VaccinationService{
			storeDB: db.StoreDB,
		},
	}
	{
		authRoute.POST("vaccination/create", vaccinationController.CreateVaccination)
		authRoute.GET("vaccination/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("vaccinations/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("vaccination/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("vaccination/:vaccination_id", vaccinationController.DeleteVaccination)
	}
}
