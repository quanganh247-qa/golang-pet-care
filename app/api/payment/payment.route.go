package payment

import (
	"net/http"

	db "github.com/quanganh247-qa/go-blog-be/app/db/sqlc"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
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
				client:  &http.Client{},
				storeDB: db.StoreDB,
			},
		},
	}

	{
		authRoute.GET("/token", goongApi.controller.GetToken)
		authRoute.GET("/banks", goongApi.controller.GetBanks)
		authRoute.POST("/generate-qr", goongApi.controller.GenerateQRCode)
	}

}
