package disease

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiceaseControllerInterface interface {
	CreateTreatment(ctx *gin.Context)
	CreateTreatmentPhase(ctx *gin.Context)
	AssignMedicineToTreatmentPhase(ctx *gin.Context)
	GetTreatmentsByPetID(ctx *gin.Context)
	GetTreatmentPhasesByTreatmentID(ctx *gin.Context)
	GetMedicinesByPhaseID(ctx *gin.Context)
	UpdateTreatmentPhaseStatus(ctx *gin.Context)
	GetActiveTreatments(ctx *gin.Context)
	GetTreatmentProgress(ctx *gin.Context)

	// getDiceaseAnhMedicinesInfo(ctx *gin.Context)
	// getTreatmentByDiseaseId(ctx *gin.Context)
}

func (c *DiceaseController) CreateTreatment(ctx *gin.Context) {
	var treatmentPhase CreateTreatmentRequest
	err := ctx.ShouldBindJSON(&treatmentPhase)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	treatment, err := c.service.CreateTreatmentService(ctx, treatmentPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment", treatment))
}

func (c *DiceaseController) CreateTreatmentPhase(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req []CreateTreatmentPhaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	treatment, err := c.service.CreateTreatmentPhaseService(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment", treatment))
}

func (c *DiceaseController) AssignMedicineToTreatmentPhase(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req []AssignMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	phase, err := c.service.AssignMedicinesToTreatmentPhase(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Phase", phase))
}

// Get Treatment By Disease ID
func (c *DiceaseController) GetTreatmentsByPetID(ctx *gin.Context) {

	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	fmt.Println("pagination", pagination, id)

	treatment, err := c.service.GetTreatmentsByPetID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment", treatment))
}

func (c *DiceaseController) GetTreatmentPhasesByTreatmentID(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	phases, err := c.service.GetTreatmentPhasesByTreatmentID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Phases", phases))
}

func (c *DiceaseController) GetMedicinesByPhaseID(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	medicines, err := c.service.GetMedicinesByPhase(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medicines", medicines))
}

func (c *DiceaseController) UpdateTreatmentPhaseStatus(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req UpdateTreatmentPhaseStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorValidator(err))
		return
	}

	err = c.service.UpdateTreatmentPhaseStatus(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Phase", nil))
}

func (c *DiceaseController) GetActiveTreatments(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	treatments, err := c.service.GetActiveTreatments(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatments", treatments))
}

func (c *DiceaseController) GetTreatmentProgress(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	progress, err := c.service.GetTreatmentProgress(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Progress", progress))
}

// func (c *DiceaseController) getDiceaseAnhMedicinesInfo(ctx *gin.Context) {
// 	disease := ctx.Query("disease")
// 	info, err := c.service.GetDiceaseAnhMedicinesInfoService(ctx, disease)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Information", info))

// }

// func (c *DiceaseController) getTreatmentByDiseaseId(ctx *gin.Context) {
// 	diseaseID := ctx.Param("disease_id")
// 	id, err := strconv.ParseInt(diseaseID, 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	treatment, err := c.service.GetTreatmentByDiseaseID(ctx, id, pagination)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment", treatment))
// }
