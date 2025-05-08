package disease

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/models"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiseaseControllerInterface interface {
	CreateDisease(ctx *gin.Context)
	CreateTreatment(ctx *gin.Context)
	CreateTreatmentPhase(ctx *gin.Context)
	AssignMedicineToTreatmentPhase(ctx *gin.Context)
	GetTreatmentsByPetID(ctx *gin.Context)
	GetTreatmentPhasesByTreatmentID(ctx *gin.Context)
	GetMedicinesByPhaseID(ctx *gin.Context)
	UpdateTreatmentPhaseStatus(ctx *gin.Context)
	GetActiveTreatments(ctx *gin.Context)
	GetTreatmentProgress(ctx *gin.Context)
	UpdateTreatmentStatus(ctx *gin.Context)
	GenerateMedicineOnlyPrescription(ctx *gin.Context)

	CreateAllergy(ctx *gin.Context)
	GetAllergiesByPetID(ctx *gin.Context)
	// GetAllMedicines(ctx *gin.Context) // Add this new method
	AnalyzeSymptoms(ctx *gin.Context)
}

// @Summary Create a new disease
// @Description Create a new disease with the given details
// @Tags disease
// @Accept json
// @Produce json
// @Param disease body CreateDiseaseRequest true "Disease details"
// @Success 200 {object} SuccessResponse "Disease created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request data"
// @Failure 500 {object} ErrorResponse "Failed to create disease"
// @Router /disease [post]
func (c *DiseaseController) CreateDisease(ctx *gin.Context) {
	var disease CreateDiseaseRequest
	if err := ctx.ShouldBindJSON(&disease); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	res, err := c.service.CreateDisease(ctx, disease)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create disease",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Disease created successfully",
		"data":    res,
	})
}

func (c *DiseaseController) CreateTreatment(ctx *gin.Context) {
	var treatmentPhase CreateTreatmentRequest
	if err := ctx.ShouldBindJSON(&treatmentPhase); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment data",
			"error":   err.Error(),
		})
		return
	}
	treatment, err := c.service.CreateTreatmentService(ctx, treatmentPhase)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create treatment",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatment created successfully",
		"data":    treatment,
	})
}

func (c *DiseaseController) CreateTreatmentPhase(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment ID",
			"error":   err.Error(),
		})
		return
	}

	var req []CreateTreatmentPhaseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid phase data",
			"error":   err.Error(),
		})
		return
	}

	treatment, err := c.service.CreateTreatmentPhaseService(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create treatment phase",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatment phase created successfully",
		"data":    treatment,
	})
}

func (c *DiseaseController) AssignMedicineToTreatmentPhase(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid phase ID",
			"error":   err.Error(),
		})
		return
	}

	var req []AssignMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid medicine assignment data",
			"error":   err.Error(),
		})
		return
	}

	phase, err := c.service.AssignMedicinesToTreatmentPhase(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to assign medicines",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Medicines assigned successfully",
		"data":    phase,
	})
}

func (c *DiseaseController) GetTreatmentsByPetID(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pagination parameters",
			"error":   err.Error(),
		})
		return
	}

	treatment, err := c.service.GetTreatmentsByPetID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch treatments",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatments retrieved successfully",
		"data":    treatment,
	})
}

func (c *DiseaseController) GetTreatmentPhasesByTreatmentID(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment ID",
			"error":   err.Error(),
		})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pagination parameters",
			"error":   err.Error(),
		})
		return
	}

	phases, err := c.service.GetTreatmentPhasesByTreatmentID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch treatment phases",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatment phases retrieved successfully",
		"data":    phases,
	})
}

func (c *DiseaseController) GetMedicinesByPhaseID(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid phase ID",
			"error":   err.Error(),
		})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pagination parameters",
			"error":   err.Error(),
		})
		return
	}

	medicines, err := c.service.GetMedicinesByPhase(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch medicines",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Medicines retrieved successfully",
		"data":    medicines,
	})
}

func (c *DiseaseController) UpdateTreatmentPhaseStatus(ctx *gin.Context) {
	phaseID := ctx.Param("phase_id")
	id, err := strconv.ParseInt(phaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid phase ID",
			"error":   err.Error(),
		})
		return
	}

	var req UpdateTreatmentPhaseStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status update data",
			"error":   err.Error(),
		})
		return
	}

	err = c.service.UpdateTreatmentPhaseStatus(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update phase status",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Phase status updated successfully",
	})
}

func (c *DiseaseController) GetActiveTreatments(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pagination parameters",
			"error":   err.Error(),
		})
		return
	}

	treatments, err := c.service.GetActiveTreatments(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch active treatments",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Active treatments retrieved successfully",
		"data":    treatments,
	})
}

func (c *DiseaseController) GetTreatmentProgress(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment ID",
			"error":   err.Error(),
		})
		return
	}

	progress, err := c.service.GetTreatmentProgress(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch treatment progress",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatment progress retrieved successfully",
		"data":    progress,
	})
}

func (c *DiseaseController) GenerateMedicineOnlyPrescription(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment ID",
			"error":   err.Error(),
		})
		return
	}

	prescription, err := c.service.GenerateMedicineOnlyPrescriptionPDF(ctx, id, "prescription.pdf")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate prescription",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Prescription generated successfully",
		"data":    prescription,
	})
}

func (c *DiseaseController) CreateAllergy(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
		return
	}

	var req CreateAllergyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid allergy data",
			"error":   err.Error(),
		})
		return
	}

	allergy, err := c.service.CreateAllergyService(ctx, id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create allergy",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Allergy created successfully",
		"data":    allergy,
	})
}

func (c *DiseaseController) GetAllergiesByPetID(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pet ID",
			"error":   err.Error(),
		})
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid pagination parameters",
			"error":   err.Error(),
		})
		return
	}
	allergies, err := c.service.GetAllergiesByPetID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch allergies",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Allergies retrieved successfully",
		"data":    allergies,
	})
}

// // Add this method to your DiseaseController implementation
// func (c *DiseaseController) GetAllMedicines(ctx *gin.Context) {
// 	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": "Invalid pagination parameters",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	medicines, err := c.service.GetAllMedicinesService(ctx, pagination)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"message": "Failed to fetch medicines",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "Medicines retrieved successfully",
// 		"data":    medicines,
// 	})
// }

// AnalyzeSymptoms processes symptom data through AI and returns treatment recommendations
func (h *DiseaseController) AnalyzeSymptoms(c *gin.Context) {
	var request SymptomAnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if llmService's client is initialized
	if h.llmService.IsInitialized() == false {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "AI service is not available. Please check API key configuration.",
		})
		return
	}

	// Validate that the patient exists
	patient, err := db.StoreDB.GetPetByID(c, request.PatientID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Patient not found"})
		return
	}

	patientModel := models.Patient{
		ID:     patient.Petid,
		Name:   patient.Name,
		Breed:  patient.Breed.String,
		Age:    float64(patient.Age.Int32),
		Gender: patient.Gender.String,
		Weight: float64(patient.Weight.Float64),
	}

	// Process the symptoms through the AI service
	treatmentPlan, err := h.llmService.AnalyzeSymptoms(patientModel, request.Symptoms)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze symptoms: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, treatmentPlan)
}

func (c *DiseaseController) UpdateTreatmentStatus(ctx *gin.Context) {
	treatmentID := ctx.Param("treatment_id")
	id, err := strconv.ParseInt(treatmentID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid treatment ID",
			"error":   err.Error(),
		})
		return
	}

	var req UpdateTreatmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid status update data",
			"error":   err.Error(),
		})
		return
	}

	err = c.service.UpdateTreatmentStatus(ctx, id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update treatment status",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Treatment status updated successfully",
	})
}
