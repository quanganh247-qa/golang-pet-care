package clinic_reporting

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func Routes(routerGroup middleware.RouterGroup) {
	reporting := routerGroup.RouterDefault.Group("/clinic-reporting")
	// authRoute := routerGroup.RouterAuth(reporting)
	perRoute := routerGroup.RouterPermission(reporting)

	controller := &ClinicReportingController{
		service: &ClinicReportingService{
			storeDB: db.StoreDB, // Replace with the actual store instance if needed
		},
	}

	// Example: Only allow users with the "ViewReports" permission to access the endpoint.
	perRoute([]perms.Permission{perms.ViewReports}).GET("", controller.GetClinicAnalytics)
}
