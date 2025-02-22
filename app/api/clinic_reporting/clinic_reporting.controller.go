package clinic_reporting

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

// ClinicReportingControllerInterface defines controller methods for clinic analytics.
type ClinicReportingControllerInterface interface {
	GetClinicAnalytics(ctx *gin.Context)
}

// ClinicReportingController implements the clinic analytics controller.
type ClinicReportingController struct {
	service ClinicReportingServiceInterface
}

// GetClinicAnalytics handles GET requests to retrieve overall clinic analytics.
func (c *ClinicReportingController) GetClinicAnalytics(ctx *gin.Context) {
	report, err := c.service.GetClinicAnalytics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, util.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, util.SuccessResponse("Clinic Analytics Report", report))
}
