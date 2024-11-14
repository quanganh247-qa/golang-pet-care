package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/handlers"
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/doctor"
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
	"github.com/quanganh247-qa/go-blog-be/app/api/medical_records"
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
>>>>>>> 0fb3f30 (user images)
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 6c35562 (dicease and treatment plan)
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/products"
	"github.com/quanganh247-qa/go-blog-be/app/api/rooms"
	"github.com/quanganh247-qa/go-blog-be/app/api/search"
=======
>>>>>>> e01abc5 (pet schedule api)
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/api/vaccination"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config, es *elasticsearch.ESService) {
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")
	routerDefault.Use(middleware.CORSMiddleware())
	routerDefault.Use(middleware.LoggingMiddleware())

	// Setup route handlers
	chatHandler := handlers.NewChatHandler(config.GoogleAPIKey, config.OpenFDAAPIKey)

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()

	v1 := routerDefault.Group(util.Configs.ApiPrefix)
	router := v1.Group("/")
	routerGroup := middleware.RouterGroup{
		RouterDefault: router,
	}
<<<<<<< HEAD
	router.GET("/health", server.healthCheck)

	chatbot.Routes(routerGroup, chatHandler)
	user.Routes(routerGroup, taskDistributor, config)
=======
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	user.Routes(routerGroup)
	service_type.Routes(routerGroup)
>>>>>>> 9d28896 (image pet)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
<<<<<<< HEAD
	appointment.Routes(routerGroup, taskDistributor)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup, es)
	petschedule.Routes(routerGroup, &config)
	vaccination.Routes(routerGroup)
	location.Routes(routerGroup, &config)
	payment.Routes(routerGroup, &config)
	cart.Routes(routerGroup)
	products.Routes(routerGroup)
	medical_records.Routes(routerGroup)
	search.Routes(routerGroup, es)
	medications.Routes(routerGroup, es)
	doctor.Routes(routerGroup)
	rooms.Routes(routerGroup)
=======
	appointment.Routes(routerGroup)
	// medications.Routes(routerGroup)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup)
	petschedule.Routes(routerGroup)
	server.Router = routerDefault
>>>>>>> 79a3bcc (medicine api)

	server.Router = routerDefault
}
