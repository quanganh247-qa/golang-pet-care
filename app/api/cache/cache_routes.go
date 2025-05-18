package cache

import (
	"github.com/quanganh247-qa/go-blog-be/app/middleware"
)

// Routes registers all cache-related routes
func Routes(routerGroup middleware.RouterGroup) {
	// Apply auth middleware - only administrators should be able to clear the cache
	cacheRouter := routerGroup.RouterDefault.Group("/cache")
	// Register the clear cache endpoint
	cacheRouter.POST("/clear", ClearAllCache)
}
