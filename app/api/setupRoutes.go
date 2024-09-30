package api

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/api/component"
	"github.com/quanganh247-qa/go-blog-be/app/api/page"
	"github.com/quanganh247-qa/go-blog-be/app/api/project"
	"github.com/quanganh247-qa/go-blog-be/app/api/user"
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func (server *Server) SetupRoutes() {
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
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/static/swagger.json")))
	user.Routes(routerGroup)
	project.Routes(routerGroup)
	page.Routes(routerGroup)
	component.Routes(routerGroup)
	server.Router = routerDefault

}
