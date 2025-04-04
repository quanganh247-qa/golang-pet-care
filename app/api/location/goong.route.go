package location

import (
	"net/http"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func Routes(routerGroup middleware.RouterGroup, config *util.Config) {
	Goong := routerGroup.RouterDefault.Group("/location")
	authRoute := routerGroup.RouterAuth(Goong)
	// Goong.Use(middleware.IPbasedRateLimitingMiddleware())

	// Apply cache middleware with longer duration for location data
	Goong.Use(middleware.CacheMiddleware(time.Hour*24, "location", []string{"GET"}))

	// Khoi tao api
	goongApi := &GoongApi{
		&GoongController{
			service: &GoongService{
				config: &GoongConfig{
					APIKey:  config.GoongAPIKey,
					BaseURL: config.GoongBaseURL,
				},
				client: &http.Client{},
			},
		},
	}

	{
		authRoute.GET("/places/autocomplete", goongApi.controller.Autocomplete)
		authRoute.GET("/places/detail", goongApi.controller.GetPlaceDetail)
		authRoute.GET("/directions", goongApi.controller.GetDirection)
		authRoute.GET("/geocode/forward", goongApi.controller.ForwardGeocode)
		authRoute.GET("/geocode/reverse", goongApi.controller.ReverseGeocode)
		// authRoute.GET("/distance", goongApi.controller.DistanceMatrix)
	}

}
