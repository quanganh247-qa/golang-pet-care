package medical_records

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicalRecordControllerInterface interface {
	CreateMedicalRecord(ctx *gin.Context)
	CreateMedicalHistory(ctx *gin.Context)
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
	diseaseID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	medicalRecord, err := db.StoreDB.GetMedicalRecord(ctx, id)
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
