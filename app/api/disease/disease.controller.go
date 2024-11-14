package disease

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

type DiceaseControllerInterface interface {
	getDiceaseAnhMedicinesInfo(ctx *gin.Context)
	getDiseaseTreatmentPlanWithPhases(ctx *gin.Context)
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

func (c *DiceaseController) getDiseaseTreatmentPlanWithPhases(ctx *gin.Context) {
	disease := ctx.Query("disease")
	treatment, err := c.service.GetDiseaseTreatmentPlanWithPhasesService(ctx, disease)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Treatment Plan", treatment))

}
