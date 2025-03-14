package api

import (
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"io"
	"net/http"
	"os"

>>>>>>> 98e9e45 (ratelimit and recovery function)
=======
>>>>>>> ada3717 (Docker file)
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot"
	"github.com/quanganh247-qa/go-blog-be/app/api/chatbot/handlers"
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
>>>>>>> c449ffc (feat: cart api)
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/clinic_reporting"
>>>>>>> ffc9071 (AI suggestion)
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/doctor"
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
	"github.com/quanganh247-qa/go-blog-be/app/api/medical_records"
<<<<<<< HEAD
<<<<<<< HEAD
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
>>>>>>> 0fb3f30 (user images)
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
<<<<<<< HEAD
=======
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/api/notification"
>>>>>>> 3bf345d (happy new year)
=======
>>>>>>> ada3717 (Docker file)
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
>>>>>>> 79a3bcc (medicine api)
=======
>>>>>>> 6c35562 (dicease and treatment plan)
=======
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
>>>>>>> 4625843 (added goong maps api)
	"github.com/quanganh247-qa/go-blog-be/app/api/notification"
<<<<<<< HEAD
>>>>>>> 9fd7fc8 (feat: validate notification schema and APIs)
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
>>>>>>> c449ffc (feat: cart api)
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/products"
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/api/rooms"
	"github.com/quanganh247-qa/go-blog-be/app/api/search"
=======
>>>>>>> e01abc5 (pet schedule api)
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/products"
>>>>>>> bd5945b (get list products)
=======
	"github.com/quanganh247-qa/go-blog-be/app/api/search"
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/api/vaccination"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
=======
>>>>>>> 6610455 (feat: redis queue)
=======
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
>>>>>>> e859654 (Elastic search)
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config, es *elasticsearch.ESService) {
=======
func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor) {
>>>>>>> 6610455 (feat: redis queue)
=======
func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config) {
>>>>>>> 1a9e82a (reset password api)
=======
func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config, es *elasticsearch.ESService) {
>>>>>>> e859654 (Elastic search)
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")
	routerDefault.Use(middleware.CORSMiddleware())
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	routerDefault.Use(middleware.LoggingMiddleware())
=======
	routerDefault.Use(middleware.IPbasedRateLimitingMiddleware())
=======
	// routerDefault.Use(middleware.IPbasedRateLimitingMiddleware())
>>>>>>> 9ee4f0a (fix bug ratelimit)
=======
	routerDefault.Use(middleware.LoggingMiddleware())
	// routerDefault.Use(middleware.ContentSecurityPolicyMiddleware())

>>>>>>> ada3717 (Docker file)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// debug := true
	// routerDefault.Use(util.Recover(logger, debug))

<<<<<<< HEAD
	// Create a custom logger with the desired output format
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr)
>>>>>>> 98e9e45 (ratelimit and recovery function)

	// Setup route handlers
	chatHandler := handlers.NewChatHandler(config.GoogleAPIKey, config.OpenFDAAPIKey)

	// logger, _ := zap.NewProduction()
	// defer logger.Sync()
=======
	// gin.DefaultWriter = io.MultiWriter(os.Stdout)
	// gin.DefaultErrorWriter = io.MultiWriter(os.Stderr)
>>>>>>> ada3717 (Docker file)

	v1 := routerDefault.Group(util.Configs.ApiPrefix)
	router := v1.Group("/")
	routerGroup := middleware.RouterGroup{
		RouterDefault: router,
	}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	router.GET("/health", server.healthCheck)

	chatbot.Routes(routerGroup, chatHandler)
	user.Routes(routerGroup, taskDistributor, config)
=======
=======

>>>>>>> 7a9ad08 (updated pet api)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
=======
	router.GET("/health", server.healthCheck)

	// // Adding the SuperTokens middleware
	// router.Use(func(c *gin.Context) {
	// 	supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// 		c.Next()
	// 	})).ServeHTTP(c.Writer, c.Request)
	// 	// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
	// 	c.Abort()
	// })
>>>>>>> ada3717 (Docker file)

	user.Routes(routerGroup, taskDistributor, config)
<<<<<<< HEAD
	service_type.Routes(routerGroup)
>>>>>>> 9d28896 (image pet)
=======
>>>>>>> b393bb9 (add service and add permission)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
<<<<<<< HEAD
<<<<<<< HEAD
	appointment.Routes(routerGroup, taskDistributor)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup, es)
	petschedule.Routes(routerGroup, &config)
<<<<<<< HEAD
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
=======
	appointment.Routes(routerGroup, taskDistributor)
>>>>>>> e859654 (Elastic search)
	device_token.Routes(routerGroup)
	disease.Routes(routerGroup, es)
	petschedule.Routes(routerGroup, &config)
	notification.Routes(routerGroup)
=======
>>>>>>> ada3717 (Docker file)
	vaccination.Routes(routerGroup)
	location.Routes(routerGroup, &config)
	payment.Routes(routerGroup, &config)
	cart.Routes(routerGroup)
	products.Routes(routerGroup)
	medical_records.Routes(routerGroup)
	clinic_reporting.Routes(routerGroup)
	search.Routes(routerGroup, es)
	medications.Routes(routerGroup, es)
	server.Router = routerDefault
>>>>>>> 79a3bcc (medicine api)

	server.Router = routerDefault
}
