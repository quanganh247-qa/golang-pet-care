package medical_records

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
=======
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicalRecordControllerInterface interface {
	CreateMedicalRecord(ctx *gin.Context)
	CreateMedicalHistory(ctx *gin.Context)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	ListMedicalHistory(ctx *gin.Context)
	GetMedicalRecord(ctx *gin.Context)
	GetMedicalHistoryByID(ctx *gin.Context)
=======
>>>>>>> 3bf345d (happy new year)
=======
	ListMedicalHistory(ctx *gin.Context)
	GetMedicalRecord(ctx *gin.Context)
	GetMedicalHistoryByID(ctx *gin.Context)
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 3bf345d (happy new year)
=======
	ListMedicalHistory(ctx *gin.Context)
	GetMedicalRecord(ctx *gin.Context)
	GetMedicalHistoryByID(ctx *gin.Context)
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
=======
	diseaseID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
>>>>>>> 3bf345d (happy new year)
=======
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
>>>>>>> e859654 (Elastic search)
=======
	diseaseID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(diseaseID, 10, 64)
>>>>>>> 3bf345d (happy new year)
=======
	petID := ctx.Param("pet_id")
	id, err := strconv.ParseInt(petID, 10, 64)
>>>>>>> e859654 (Elastic search)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	medicalRecord, err := c.service.GetMedicalRecord(ctx, id)
=======
	medicalRecord, err := db.StoreDB.GetMedicalRecord(ctx, id)
>>>>>>> 3bf345d (happy new year)
=======
	medicalRecord, err := c.service.GetMedicalRecord(ctx, id)
>>>>>>> e859654 (Elastic search)
=======
	medicalRecord, err := db.StoreDB.GetMedicalRecord(ctx, id)
>>>>>>> 3bf345d (happy new year)
=======
	medicalRecord, err := c.service.GetMedicalRecord(ctx, id)
>>>>>>> e859654 (Elastic search)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> e859654 (Elastic search)

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

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> e859654 (Elastic search)
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

<<<<<<< HEAD
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 4ccd381 (Update appointment flow)
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 4ccd381 (Update appointment flow)
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
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
=======
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> e859654 (Elastic search)
