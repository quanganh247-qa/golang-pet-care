package payment

import (
	"net/http"

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
=======
>>>>>>> c449ffc (feat: cart api)
=======
	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
>>>>>>> b0fe977 (place order and make payment)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	payment := routerGroup.RouterDefault.Group("/payment")
	authRoute := routerGroup.RouterAuth(payment)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 6b24d88 (feat(payment): add PayOS payment integration and enhance treatment module)
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
=======
	payment := routerGroup.RouterDefault.Group("/payment")
	authRoute := routerGroup.RouterAuth(payment)
>>>>>>> e859654 (Elastic search)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
					// PaymentAPIKey:    config.PaymentAPIKey,
					// PaymentBaseURL:   config.PaymentBaseURL,
					// PaymentClientKey: config.PaymentClientKey,
=======
					PayPalClientID:     config.PaypalClientID,
					PayPalClientSecret: config.PaypalClientSecret,
					PayPalBaseURL:      config.PaypalURL,
>>>>>>> ada3717 (Docker file)
				},
<<<<<<< HEAD
				client: &http.Client{},
>>>>>>> c449ffc (feat: cart api)
=======
				client:  &http.Client{},
				storeDB: db.StoreDB,
>>>>>>> b0fe977 (place order and make payment)
=======
	Goong := routerGroup.RouterDefault.Group("/vietqr")
	authRoute := routerGroup.RouterAuth(Goong)
=======
	payment := routerGroup.RouterDefault.Group("/payment")
	authRoute := routerGroup.RouterAuth(payment)
>>>>>>> e859654 (Elastic search)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Khoi tao api
	paymentApi := &PaymentApi{
		&PaymentController{
			service: &PaymentService{
				config: &PaymentConfig{
					// PaymentAPIKey:    config.PaymentAPIKey,
					// PaymentBaseURL:   config.PaymentBaseURL,
					// PaymentClientKey: config.PaymentClientKey,
=======
					PayPalClientID:     config.PaypalClientID,
					PayPalClientSecret: config.PaypalClientSecret,
					PayPalBaseURL:      config.PaypalURL,
>>>>>>> ada3717 (Docker file)
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
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		authRoute.GET("/token", paymentApi.controller.GetToken)
		authRoute.GET("/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
=======
		authRoute.GET("/token", goongApi.controller.GetToken)
		authRoute.GET("/banks", goongApi.controller.GetBanks)
		authRoute.POST("/generate-qr", goongApi.controller.GenerateQRCode)
>>>>>>> c449ffc (feat: cart api)
=======
		authRoute.GET("/token", paymentApi.controller.GetToken)
		authRoute.GET("/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
>>>>>>> e859654 (Elastic search)
=======
		authRoute.GET("/token", goongApi.controller.GetToken)
		authRoute.GET("/banks", goongApi.controller.GetBanks)
		authRoute.POST("/generate-qr", goongApi.controller.GenerateQRCode)
>>>>>>> c449ffc (feat: cart api)
=======
		authRoute.GET("/token", paymentApi.controller.GetToken)
		authRoute.GET("/banks", paymentApi.controller.GetBanks)
		authRoute.POST("/generate-qr", paymentApi.controller.GenerateQRCode)
>>>>>>> e859654 (Elastic search)
	}

}
