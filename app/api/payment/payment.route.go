package payment

import (
	"net/http"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
	payment := routerGroup.RouterDefault.Group("/")
	authRoute := routerGroup.RouterAuth(payment)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
					VietQRAPIKey:    config.VietQRAPIKey,
					VietQRClientKey: config.VietQRClientKey,
					VietQRBaseURL:   config.VietQRBaseURL,
				},
				client:  &http.Client{},
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.GET("/payment/token", paymentApi.controller.GetToken)
		authRoute.GET("/payment/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/payment/generate-qr", paymentApi.controller.GenerateQRCode)
		authRoute.POST("/payment/quick-link", paymentApi.controller.GenerateQuickLink)
	}
	{
		// Add this route to your existing routes
		authRoute.GET("/payment/revenue/last-seven-days", paymentApi.controller.GetRevenueLastSevenDays)
		// Add the new patient trends endpoint
		authRoute.GET("/patients/trends", paymentApi.controller.GetPatientTrends)
	}
}
