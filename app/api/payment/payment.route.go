package payment

import (
	"net/http"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
	payment := routerGroup.RouterDefault.Group("/payment")
	authRoute := routerGroup.RouterAuth(payment)

	// Initialize API
	paymentApi := &PaymentApi{
		controller: &PaymentController{
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

	// Existing routes
	authRoute.GET("/token", paymentApi.controller.GetToken)
	authRoute.GET("/banks", paymentApi.controller.GetBanks)
	authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
	authRoute.POST("/quick-link", paymentApi.controller.GenerateQuickLink)
	authRoute.GET("/revenue/last-seven-days", paymentApi.controller.GetRevenueLastSevenDays)
	authRoute.POST("/cash", paymentApi.controller.CreateCashPayment)
	authRoute.POST("/confirm", paymentApi.controller.ConfirmPayment)
	// Add new route for listing payments
	authRoute.GET("/list", paymentApi.controller.ListPayments)
}
