package smtp

import (
	"github.com/gin-gonic/gin"
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

// RegisterRoutes registers all SMTP configuration routes
func RegisterRoutes(router *gin.RouterGroup, config util.Config, store db.Store) {
	smtpService := NewSMTPConfigService(config, store)

	// Create router group
	smtpRoutes := router.Group("/smtp")

	// Create middleware router group
	routerGroup := middleware.RouterGroup{
		RouterDefault: smtpRoutes,
	}

	// Sử dụng permission middleware cho admin
	perRoute := routerGroup.RouterPermission(smtpRoutes)

	// Đăng ký routes với permission
	perRoute([]perms.Permission{perms.ManageSystem}).POST("/config", smtpService.CreateSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).GET("/config/:id", smtpService.GetSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).GET("/config", smtpService.ListSMTPConfigs)
	perRoute([]perms.Permission{perms.ManageSystem}).GET("/config/default", smtpService.GetDefaultSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).PUT("/config/:id", smtpService.UpdateSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).DELETE("/config/:id", smtpService.DeleteSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).POST("/config/default", smtpService.SetDefaultSMTPConfig)
	perRoute([]perms.Permission{perms.ManageSystem}).POST("/config/test", smtpService.TestSMTPConfig)
}
