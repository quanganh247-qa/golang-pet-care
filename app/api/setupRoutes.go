package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/handlers"
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
	"github.com/quanganh247-qa/go-blog-be/app/api/doctor"
	"github.com/quanganh247-qa/go-blog-be/app/api/invoice"
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
	"github.com/quanganh247-qa/go-blog-be/app/api/medical_records"
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
	"github.com/quanganh247-qa/go-blog-be/app/api/products"
	"github.com/quanganh247-qa/go-blog-be/app/api/reports"
	"github.com/quanganh247-qa/go-blog-be/app/api/rooms"
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/smtp"
	"github.com/quanganh247-qa/go-blog-be/app/api/test"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/api/vaccination"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/inference"
	"github.com/quanganh247-qa/go-blog-be/app/service/minio"
	"github.com/quanganh247-qa/go-blog-be/app/service/websocket"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config, ws *websocket.WSClientManager) {
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")

	// Apply global security middlewares
	routerDefault.Use(middleware.LoggingMiddleware())
	routerDefault.Use(middleware.CORSMiddleware())

	// Setup route handlers
	chatHandler := handlers.NewChatHandler(config.GoogleAPIKey)
	// Add Roboflow inference handler
	inferenceHandler := inference.NewInferenceHandler(config.RoboflowAPIKey)

	v1 := routerDefault.Group(util.Configs.ApiPrefix)
	router := v1.Group("/")
	routerGroup := middleware.RouterGroup{
		RouterDefault: router,
	}
	routerDefault.GET("/ws", server.ws.HandleWebSocket)

	// Health check endpoint
	router.GET("/health", server.healthCheck)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register all API routes
	chatbot.Routes(routerGroup, chatHandler)
	// Register inference routes
	inferenceHandler.RegisterRoutes(router)
	user.Routes(routerGroup, taskDistributor, config)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
	appointment.Routes(routerGroup, taskDistributor, ws)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup)
	petschedule.Routes(routerGroup, &config)
	vaccination.Routes(routerGroup)
	location.Routes(routerGroup, &config)
	payment.Routes(routerGroup, &config)
	cart.Routes(routerGroup)
	products.Routes(routerGroup)
	medical_records.Routes(routerGroup)
	test.Routes(routerGroup, ws)
	medications.Routes(routerGroup, taskDistributor, ws)
	doctor.Routes(routerGroup)
	rooms.Routes(routerGroup)
	invoice.Routes(routerGroup)
	reports.Routes(routerGroup)
	// Register SMTP configuration routes
	smtp.RegisterRoutes(router, config, server.store)

	minioHandler := minio.NewMinioHandler(routerDefault)
	minio.Routes(routerGroup, minioHandler)
	server.Router = routerDefault
}
