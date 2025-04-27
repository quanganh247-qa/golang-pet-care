package medical_records

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicalRecordControllerInterface interface {
	CreateMedicalRecord(ctx *gin.Context)
	CreateMedicalHistory(ctx *gin.Context)
	ListMedicalHistory(ctx *gin.Context)
	GetMedicalRecord(ctx *gin.Context)
	GetMedicalHistoryByID(ctx *gin.Context)
	GetMedicalHistoryByPetID(ctx *gin.Context)

	// New methods for expanded functionality
	CreateExamination(ctx *gin.Context)
	GetExamination(ctx *gin.Context)
	ListExaminations(ctx *gin.Context)

	CreatePrescription(ctx *gin.Context)
	GetPrescription(ctx *gin.Context)
	ListPrescriptions(ctx *gin.Context)

	CreateTestResult(ctx *gin.Context)
	GetTestResult(ctx *gin.Context)
	ListTestResults(ctx *gin.Context)

	GetCompleteMedicalHistory(ctx *gin.Context)
}

func (c *MedicalRecordController) CreateMedicalRecord(ctx *gin.Context) {
	diseaseID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.CreateMedicalRecord(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical Record", res))

}

func (c *MedicalRecordController) CreateMedicalHistory(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	medicalRecord, err := c.service.GetMedicalRecord(ctx, id)
	if err != nil {
		medicalRecord, err = c.service.CreateMedicalRecord(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	var req MedicalHistoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.CreateMedicalHistory(ctx, &req, medicalRecord.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical History", res))
}

func (c *MedicalRecordController) ListMedicalHistory(ctx *gin.Context) {
	medicalRecordID := ctx.Param("medical_record_id")
	id, err := strconv.ParseInt(medicalRecordID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.ListMedicalHistory(ctx, id, pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical History", res))
}

func (c *MedicalRecordController) GetMedicalRecord(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetMedicalRecord(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical Record", res))
}

func (c *MedicalRecordController) GetMedicalHistoryByID(ctx *gin.Context) {
	medicalHistoryID := ctx.Param("id")
	id, err := strconv.ParseInt(medicalHistoryID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetMedicalHistoryByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical History", res))
}

func (c *MedicalRecordController) GetMedicalHistoryByPetID(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	res, err := c.service.GetMedicalHistoryByPetID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical History", res))
}

// New controller methods for examinations
func (c *MedicalRecordController) CreateExamination(ctx *gin.Context) {
	var req ExaminationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.CreateExamination(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Examination", res))
}

func (c *MedicalRecordController) GetExamination(ctx *gin.Context) {
	examinationID := ctx.Param("examination_id")
	id, err := strconv.ParseInt(examinationID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetExamination(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Examination", res))
}

func (c *MedicalRecordController) ListExaminations(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	recordID := ctx.Query("medical_record_id")

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var res []ExaminationResponse

	if petID != "" {
		// List by pet ID
		id, err := strconv.ParseInt(petID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListExaminationsByPet(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if recordID != "" {
		// List by medical record ID
		id, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListExaminationsByMedicalHistory(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either pet_id or medical_record_id is required"})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Examinations", res))
}

// New controller methods for prescriptions
func (c *MedicalRecordController) CreatePrescription(ctx *gin.Context) {
	var req PrescriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.CreatePrescription(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Prescription", res))
}

func (c *MedicalRecordController) GetPrescription(ctx *gin.Context) {
	prescriptionID := ctx.Param("prescription_id")
	id, err := strconv.ParseInt(prescriptionID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetPrescription(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Prescription", res))
}

func (c *MedicalRecordController) ListPrescriptions(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	recordID := ctx.Query("medical_record_id")

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var res []PrescriptionResponse

	if petID != "" {
		// List by pet ID
		id, err := strconv.ParseInt(petID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListPrescriptionsByPet(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if recordID != "" {
		// List by medical record ID
		id, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListPrescriptionsByMedicalHistory(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either pet_id or medical_record_id is required"})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Prescriptions", res))
}

// New controller methods for test results
func (c *MedicalRecordController) CreateTestResult(ctx *gin.Context) {
	var req TestResultRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.CreateTestResult(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Test Result", res))
}

func (c *MedicalRecordController) GetTestResult(ctx *gin.Context) {
	testResultID := ctx.Param("test_result_id")
	id, err := strconv.ParseInt(testResultID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetTestResult(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Test Result", res))
}

func (c *MedicalRecordController) ListTestResults(ctx *gin.Context) {
	petID := ctx.Query("pet_id")
	recordID := ctx.Query("medical_record_id")

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var res []TestResultResponse

	if petID != "" {
		// List by pet ID
		id, err := strconv.ParseInt(petID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListTestResultsByPet(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if recordID != "" {
		// List by medical record ID
		id, err := strconv.ParseInt(recordID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
			return
		}

		res, err = c.service.ListTestResultsByMedicalHistory(ctx, id, pagination)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Either pet_id or medical_record_id is required"})
		return
	}

	ctx.JSON(http.StatusOK, util.SuccessResponse("Test Results", res))
}

// Complete medical history method
func (c *MedicalRecordController) GetCompleteMedicalHistory(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.GetCompleteMedicalHistory(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Medical History", res))
}
