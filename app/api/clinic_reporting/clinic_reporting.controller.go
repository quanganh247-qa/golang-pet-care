package clinic_reporting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// ClinicReportingControllerInterface defines controller methods for clinic analytics.
type ClinicReportingControllerInterface interface {
	GenerateSOAPNote(ctx *gin.Context)
}

// ClinicReportingController implements the clinic analytics controller.
type ClinicReportingController struct {
	service ClinicReportingServiceInterface
}

func (c *ClinicReportingController) GenerateSOAPNote(ctx *gin.Context) {
	var req GenerateSOAPNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, util.ErrorResponse(err))
		return
	}
	note, err := c.service.GenerateSOAPNote(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("SOAP Note", note))
}
