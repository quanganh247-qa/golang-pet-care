package vaccination

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {

<<<<<<< HEAD
	vaccination := routerGroup.RouterDefault.Group("/")
=======
	vaccination := routerGroup.RouterDefault.Group("/vaccination")
>>>>>>> 290baeb (fixed vaccine routes)
	authRoute := routerGroup.RouterAuth(vaccination)

	vaccinationController := &VaccinationController{
		service: &VaccinationService{
			storeDB: db.StoreDB,
		},
	}
	{
<<<<<<< HEAD
		authRoute.POST("vaccination/create", vaccinationController.CreateVaccination)
		authRoute.GET("vaccination/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("vaccinations/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("vaccination/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("vaccination/:vaccination_id", vaccinationController.DeleteVaccination)
=======
		authRoute.POST("/create", vaccinationController.CreateVaccination)
		authRoute.GET("/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("/:vaccination_id", vaccinationController.DeleteVaccination)
>>>>>>> 290baeb (fixed vaccine routes)
	}
}
