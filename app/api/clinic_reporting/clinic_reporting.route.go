package clinic_reporting

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

func Routes(routerGroup middleware.RouterGroup) {
	reporting := routerGroup.RouterDefault.Group("/clinic")
	// authRoute := routerGroup.RouterAuth(reporting)
	// perRoute := routerGroup.RouterPermission(reporting)

	controller := &ClinicReportingController{
		service: &ClinicReportingService{
			storeDB: db.StoreDB, // Replace with the actual store instance if needed
		},
	}
	{
		reporting.POST("/soap-note", controller.GenerateSOAPNote)
	}

	// Example: Only allow users with the "ViewReports" permission to access the endpoint.
}
