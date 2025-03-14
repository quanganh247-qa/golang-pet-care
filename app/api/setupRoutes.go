package api

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/appointment"
	"github.com/quanganh247-qa/go-blog-be/app/api/cart"
	"github.com/quanganh247-qa/go-blog-be/app/api/clinic_reporting"
	"github.com/quanganh247-qa/go-blog-be/app/api/device_token"
	"github.com/quanganh247-qa/go-blog-be/app/api/disease"
	"github.com/quanganh247-qa/go-blog-be/app/api/location"
	"github.com/quanganh247-qa/go-blog-be/app/api/medical_records"
	"github.com/quanganh247-qa/go-blog-be/app/api/medications"
	"github.com/quanganh247-qa/go-blog-be/app/api/payment"
	"github.com/quanganh247-qa/go-blog-be/app/api/pet"
	petschedule "github.com/quanganh247-qa/go-blog-be/app/api/pet_schedule"
	"github.com/quanganh247-qa/go-blog-be/app/api/products"
	"github.com/quanganh247-qa/go-blog-be/app/api/search"
	"github.com/quanganh247-qa/go-blog-be/app/api/service"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/api/vaccination"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/service/elasticsearch"
	"github.com/quanganh247-qa/go-blog-be/app/service/worker"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"go.uber.org/zap"
)

func (server *Server) SetupRoutes(taskDistributor worker.TaskDistributor, config util.Config, es *elasticsearch.ESService) {
	gin.SetMode(gin.ReleaseMode)
	routerDefault := gin.New()
	routerDefault.SetTrustedProxies(nil)
	routerDefault.Static("/static", "app/static")
	routerDefault.Use(middleware.CORSMiddleware())
	routerDefault.Use(middleware.LoggingMiddleware())
	// routerDefault.Use(middleware.ContentSecurityPolicyMiddleware())

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// debug := true
	// routerDefault.Use(util.Recover(logger, debug))

	// gin.DefaultWriter = io.MultiWriter(os.Stdout)
	// gin.DefaultErrorWriter = io.MultiWriter(os.Stderr)

	v1 := routerDefault.Group(util.Configs.ApiPrefix)
	router := v1.Group("/")
	routerGroup := middleware.RouterGroup{
		RouterDefault: router,
	}
	router.GET("/health", server.healthCheck)

	// // Adding the SuperTokens middleware
	// router.Use(func(c *gin.Context) {
	// 	supertokens.Middleware(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// 		c.Next()
	// 	})).ServeHTTP(c.Writer, c.Request)
	// 	// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
	// 	c.Abort()
	// })

	user.Routes(routerGroup, taskDistributor, config)
	pet.Routes(routerGroup)
	service.Routes(routerGroup)
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
	clinic_reporting.Routes(routerGroup)
	search.Routes(routerGroup, es)
	medications.Routes(routerGroup, es)
	server.Router = routerDefault

}
