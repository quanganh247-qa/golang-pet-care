package vaccination

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type VaccinationController struct {
	service VaccinationServiceInterface
}

func (c *VaccinationController) CreateVaccination(ctx *gin.Context) {
	var req createVaccinationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.CreateVaccination(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *VaccinationController) GetVaccinationByID(ctx *gin.Context) {
	vaccinationID, err := strconv.ParseInt(ctx.Param("vaccination_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vaccination ID"})
		return
	}

	res, err := c.service.GetVaccinationByID(ctx, vaccinationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *VaccinationController) ListVaccinationsByPetID(ctx *gin.Context) {

	petID, err := strconv.ParseInt(ctx.Param("pet_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := c.service.ListVaccinationsByPetID(ctx, petID, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (c *VaccinationController) UpdateVaccination(ctx *gin.Context) {
	vaccinationID, err := strconv.ParseInt(ctx.Param("vaccination_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vaccination ID"})
		return
	}

	var req updateVaccinationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.VaccinationID = vaccinationID

	if err := c.service.UpdateVaccination(ctx, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Vaccination updated successfully"})
}

func (c *VaccinationController) DeleteVaccination(ctx *gin.Context) {
	vaccinationID, err := strconv.ParseInt(ctx.Param("vaccination_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vaccination ID"})
		return
	}

	if err := c.service.DeleteVaccination(ctx, vaccinationID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Vaccination deleted successfully"})
}

func (c *VaccinationController) GetUpcomingVaccinations(ctx *gin.Context) {
	// Parse pet ID if provided (optional)
	var petID int64 = 0
	petIDParam := ctx.Query("pet_id")
	if petIDParam != "" {
		var err error
		petID, err = strconv.ParseInt(petIDParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pet ID"})
			return
		}
	}

	// Parse days parameter (defaults to 30 days if not provided)
	daysAhead := 30
	daysParam := ctx.Query("days")
	if daysParam != "" {
		var err error
		daysAhead, err = strconv.Atoi(daysParam)
		if err != nil || daysAhead < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid days parameter"})
			return
		}
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	res, err := c.service.GetUpcomingVaccinations(ctx, petID, daysAhead, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}
