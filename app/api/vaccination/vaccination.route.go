package vaccination

import (
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
)

func Routes(router *gin.Engine, store db.Store) {
	vaccinationController := &VaccinationController{
		service: &VaccinationService{
			storeDB: store,
		},
	}

	vaccination := router.Group("/vaccination")
	{
		vaccination.POST("/create", vaccinationController.CreateVaccination)
		vaccination.GET("/:vaccination_id", vaccinationController.GetVaccinationByID)
		vaccination.GET("/pet/:pet_id", vaccinationController.ListVaccinationsByPetID)
		vaccination.PUT("/:vaccination_id", vaccinationController.UpdateVaccination)
		vaccination.DELETE("/:vaccination_id", vaccinationController.DeleteVaccination)
	}
}
