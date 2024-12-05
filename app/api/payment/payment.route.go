package payment

import (
	"net/http"

<<<<<<< HEAD
<<<<<<< HEAD
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
=======
>>>>>>> c449ffc (feat: cart api)
=======
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> b0fe977 (place order and make payment)
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
<<<<<<< HEAD
	payment := routerGroup.RouterDefault.Group("/payment")
	authRoute := routerGroup.RouterAuth(payment)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
					VietQRAPIKey:       config.VietQRAPIKey,
					VietQRClientKey:    config.VietQRClientKey,
					VietQRBaseURL:      config.VietQRBaseURL,
					PayPalClientID:     config.PaypalClientID,
					PayPalClientSecret: config.PaypalClientSecret,
					PayPalBaseURL:      config.PaypalURL,
				},
				client:  &http.Client{},
				storeDB: db.StoreDB,
=======
	Goong := routerGroup.RouterDefault.Group("/vietqr")
	authRoute := routerGroup.RouterAuth(Goong)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	goongApi := &VietQRApi{
		&VietQRController{
			service: &VietQRService{
				config: &VietQRConfig{
					APIKey:    config.VietQRAPIKey,
					BaseURL:   config.VietQRBaseURL,
					ClientKey: config.VietQRClientKey,
				},
<<<<<<< HEAD
				client: &http.Client{},
>>>>>>> c449ffc (feat: cart api)
=======
				client:  &http.Client{},
				storeDB: db.StoreDB,
>>>>>>> b0fe977 (place order and make payment)
			},
		},
	}

	{
<<<<<<< HEAD
		authRoute.GET("/token", paymentApi.controller.GetToken)
		authRoute.GET("/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
=======
		authRoute.GET("/token", goongApi.controller.GetToken)
		authRoute.GET("/banks", goongApi.controller.GetBanks)
		authRoute.POST("/generate-qr", goongApi.controller.GenerateQRCode)
>>>>>>> c449ffc (feat: cart api)
	}

}
