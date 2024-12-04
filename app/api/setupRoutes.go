package api

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
	"github.com/quanganh247-qa/go-blog-be/app/api/notification"
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/service_type"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/api/vaccination"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config) {
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")
	routerDefault.Use(middleware.CORSMiddleware())
	// routerDefault.Use(middleware.IPbasedRateLimitingMiddleware())
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	debug := true // or false, depending on your environment
	// Apply the custom recovery middleware
	routerDefault.Use(util.Recover(logger, debug))

	// Create a custom logger with the desired output format
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr)

	// routerDefault.Use(middleware.LoggerCustom())

	v1 := routerDefault.Group(util.Configs.ApiPrefix)
	router := v1.Group("/")
	routerGroup := middleware.RouterGroup{
		RouterDefault: router,
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	user.Routes(routerGroup, taskDistributor, config)
	service_type.Routes(routerGroup)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
	appointment.Routes(routerGroup)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup)
	petschedule.Routes(routerGroup)
	notification.Routes(routerGroup)
	vaccination.Routes(routerGroup)
	location.Routes(routerGroup, &config)
	payment.Routes(routerGroup, &config)
	cart.Routes(routerGroup)

	server.Router = routerDefault

}
