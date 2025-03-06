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
	CreateAllergy(ctx *gin.Context)
	GetMedicalHistoryByID(ctx *gin.Context)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
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

func (c *MedicalRecordController) CreateAllergy(ctx *gin.Context) {
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	var req AllergyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	res, err := c.service.CreateAllergy(ctx, req, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Allergy", res))
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
