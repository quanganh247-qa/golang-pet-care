package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
)

type RouterGroup struct {
	RouterDefault *gin.RouterGroup
}

func (routerGroup *RouterGroup) RouterAuth(router *gin.RouterGroup) gin.IRoutes {
	newRouter := router.Group("/")
	return newRouter.Use(AuthMiddleware(token.TokenMaker))
}

// func (routerGroup *RouterGroup) RouterPermission(router *gin.RouterGroup, typeApi perms.TypeApi) func([]perms.Permission) gin.IRoutes {
// 	return func(method []perms.Permission) gin.IRoutes {
// 		newRouter := router.Group("/")
// 		return newRouter.Use(AuthMiddleware(token.TokenMaker), PermissionMiddleware(method,typeApi))
// 	}
// }
