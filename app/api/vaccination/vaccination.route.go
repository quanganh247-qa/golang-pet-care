package vaccination

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	vaccination := routerGroup.RouterDefault.Group("/")
=======
	vaccination := routerGroup.RouterDefault.Group("/vaccination")
>>>>>>> 290baeb (fixed vaccine routes)
=======
	vaccination := routerGroup.RouterDefault.Group("/")
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
		authRoute.POST("vaccination/create", vaccinationController.CreateVaccination)
		authRoute.GET("vaccination/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("vaccinations/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("vaccination/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("vaccination/:vaccination_id", vaccinationController.DeleteVaccination)
<<<<<<< HEAD
=======
=======
>>>>>>> 290baeb (fixed vaccine routes)
		authRoute.POST("/create", vaccinationController.CreateVaccination)
		authRoute.GET("/:vaccination_id", vaccinationController.GetVaccinationByID)
		authRoute.GET("/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		authRoute.PUT("/:vaccination_id", vaccinationController.UpdateVaccination)
		authRoute.DELETE("/:vaccination_id", vaccinationController.DeleteVaccination)
<<<<<<< HEAD
>>>>>>> 290baeb (fixed vaccine routes)
=======
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
=======
>>>>>>> 290baeb (fixed vaccine routes)
	}
}
