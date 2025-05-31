package reports

import (
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

// Routes sets up the routes for the reports API
func Routes(routerGroup middleware.RouterGroup) {
	// Initialize handler with store
	handler := NewHandler(db.StoreDB)

	// Create reports router group
	reportsRouter := routerGroup.RouterDefault.Group("/reports")

	// Permission-based routes - only users with ViewReports permission can access
	perRoute := routerGroup.RouterPermission(reportsRouter)

	// Set up report endpoints with proper permissions
	perRoute([]perms.Permission{perms.ViewReports}).GET("/financial", handler.GetFinancialReport)
	perRoute([]perms.Permission{perms.ViewReports}).GET("/medical", handler.GetMedicalRecordsReport)

	// Doctor statistics endpoints
	// perRoute([]perms.Permission{perms.ViewReports}).GET("/doctors/stats", handler.GetAllDoctorsStats)
	perRoute([]perms.Permission{perms.ViewReports}).GET("/doctors/:id/stats", handler.GetDoctorStats)

	reportsRouter.GET("/doctors", handler.GetAllDoctorsStats)
	reportsRouter.GET("/rooms", handler.ListRoomHandler)

}
