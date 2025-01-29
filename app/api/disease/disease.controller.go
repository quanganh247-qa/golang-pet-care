package disease

import (
	"net/http"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"strconv"
=======
>>>>>>> 6c35562 (dicease and treatment plan)
=======
	"strconv"
>>>>>>> 6a85052 (get treatment by disease)
=======
>>>>>>> 6c35562 (dicease and treatment plan)
=======
	"strconv"
>>>>>>> 6a85052 (get treatment by disease)

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
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
	GenerateMedicineOnlyPrescription(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD

	CreateAllergy(ctx *gin.Context)
	GetAllergiesByPetID(ctx *gin.Context)
}

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
=======
type DiceaseControllerInterface interface {
=======
type DiseaseControllerInterface interface {
	CreateDisease(ctx *gin.Context)
>>>>>>> e859654 (Elastic search)
	CreateTreatment(ctx *gin.Context)
	CreateTreatmentPhase(ctx *gin.Context)
	AssignMedicineToTreatmentPhase(ctx *gin.Context)
	GetTreatmentsByPetID(ctx *gin.Context)
	GetTreatmentPhasesByTreatmentID(ctx *gin.Context)
	GetMedicinesByPhaseID(ctx *gin.Context)
	UpdateTreatmentPhaseStatus(ctx *gin.Context)
	GetActiveTreatments(ctx *gin.Context)
	GetTreatmentProgress(ctx *gin.Context)
=======
>>>>>>> ada3717 (Docker file)
=======

	CreateAllergy(ctx *gin.Context)
	GetAllergiesByPetID(ctx *gin.Context)
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
}

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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (c *DiceaseController) getDiseaseTreatmentPlanWithPhases(ctx *gin.Context) {
	disease := ctx.Query("disease")
	treatment, err := c.service.GetDiseaseTreatmentPlanWithPhasesService(ctx, disease)
=======
func (c *DiceaseController) UpdateTreatmentPhaseStatus(ctx *gin.Context) {
=======
func (c *DiseaseController) UpdateTreatmentPhaseStatus(ctx *gin.Context) {
>>>>>>> e859654 (Elastic search)
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
>>>>>>> 883d5b3 (update treatment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update phase status",
			"error":   err.Error(),
		})
		return
	}
<<<<<<< HEAD
<<<<<<< HEAD
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment Plan", treatment))

>>>>>>> 6c35562 (dicease and treatment plan)
}

=======
>>>>>>> 3bf345d (happy new year)
func (c *DiceaseController) getTreatmentByDiseaseId(ctx *gin.Context) {
	diseaseID := ctx.Param("disease_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
=======
	ctx.JSON(http.StatusOK, util.SuccessResponse("Phase", nil))
=======
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Phase status updated successfully",
	})
>>>>>>> c8bec46 (feat: add chatbot, room management, and pet allergy features)
}

func (c *DiseaseController) GetActiveTreatments(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
>>>>>>> 883d5b3 (update treatment)
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
=======
type DiceaseControllerInterface interface {
	CreateTreatment(ctx *gin.Context)
	CreateTreatmentPhase(ctx *gin.Context)

	getDiceaseAnhMedicinesInfo(ctx *gin.Context)
	getTreatmentByDiseaseId(ctx *gin.Context)
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

func (c *DiceaseController) getDiceaseAnhMedicinesInfo(ctx *gin.Context) {
	disease := ctx.Query("disease")
	info, err := c.service.GetDiceaseAnhMedicinesInfoService(ctx, disease)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Information", info))

}

<<<<<<< HEAD
func (c *DiceaseController) getDiseaseTreatmentPlanWithPhases(ctx *gin.Context) {
	disease := ctx.Query("disease")
	treatment, err := c.service.GetDiseaseTreatmentPlanWithPhasesService(ctx, disease)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment Plan", treatment))

>>>>>>> 6c35562 (dicease and treatment plan)
}

=======
>>>>>>> 3bf345d (happy new year)
func (c *DiceaseController) getTreatmentByDiseaseId(ctx *gin.Context) {
	diseaseID := ctx.Param("disease_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	treatment, err := c.service.GetTreatmentByDiseaseID(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment", treatment))
}
