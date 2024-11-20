package api

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
	"github.com/quanganh247-qa/go-blog-be/app/api/notification"
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/service_type"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor) {
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")
	routerDefault.Use(middleware.CORSMiddleware())
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
	user.Routes(routerGroup, taskDistributor)
	service_type.Routes(routerGroup)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
	appointment.Routes(routerGroup)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup)
	petschedule.Routes(routerGroup)
	notification.Routes(routerGroup)

	server.Router = routerDefault

}
