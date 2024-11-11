package medications

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type MedicineControllerInterface interface {
	CreateMedicine(ctx *gin.Context)
	GetMedicineByID(ctx *gin.Context)
	ListMedicines(ctx *gin.Context)
	// ListMedicinesByUsername(ctx *gin.Context)
	UpdateMedicine(ctx *gin.Context)
	// DeleteMedicine(ctx *gin.Context)
}

func (c *MedicineController) CreateMedicine(ctx *gin.Context) {
	var req createMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authPayload, err := middleware.GetAuthorizationPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := c.service.CreateMedicine(ctx, authPayload.Username, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *MedicineController) GetMedicineByID(ctx *gin.Context) {
	medicineidStr := ctx.Param("medicine_id")
	medicineid, err := strconv.ParseInt(medicineidStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
		return
	}

	res, err := c.service.GetMedicineByID(ctx, medicineid)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *MedicineController) ListMedicines(ctx *gin.Context) {

	petID := ctx.Param("pet_id")
	petIDInt, err := strconv.ParseInt(petID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Pet ID"})
		return
	}

	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}

	Medicines, err := c.service.ListMedicines(ctx, pagination, petIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Medicines)
}

func (c *MedicineController) UpdateMedicine(ctx *gin.Context) {
	medicineid, err := strconv.ParseInt(ctx.Param("medicine_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
		return
	}

	var req createMedicineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(req)

	err = c.service.UpdateMedicine(ctx, medicineid, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medicine updated successfully"})
}

// func (c *MedicineController) DeleteMedicine(ctx *gin.Context) {
// 	Medicineid, err := strconv.ParseInt(ctx.Param("Medicineid"), 10, 64)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Medicine ID"})
// 		return
// 	}

// 	err = c.service.SetMedicineInactive(ctx, Medicineid)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Medicine set to inactive successfully"})
// }

// func (c *MedicineController) ListMedicinesByUsername(ctx *gin.Context) {
// 	// username := ctx.Param("username")
// 	pagination, err := util.GetPageInQuery(ctx.Request.URL.Query())
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
// 		return
// 	}
// 	authPayload, err := middleware.GetAuthorizationPayload(ctx)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Println(authPayload.Username)
// 	Medicines, err := c.service.ListMedicinesByUsername(ctx, authPayload.Username, pagination)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, Medicines)
// }
