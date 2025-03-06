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
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
					// PaymentAPIKey:    config.PaymentAPIKey,
					// PaymentBaseURL:   config.PaymentBaseURL,
					// PaymentClientKey: config.PaymentClientKey,
				},
				client:  &http.Client{},
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.GET("/token", paymentApi.controller.GetToken)
		authRoute.GET("/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
	}

}
