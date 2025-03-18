package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigins := strings.Split(util.Configs.AccessControlAllowOrigin, ",")

		// Check if origin is in allowed list or if we're in development mode
		if origin != "" {
			// Either allow specific origins from config or allow the requesting origin
			if slices.Contains(allowedOrigins, origin) || strings.Contains(origin, "localhost") {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else if len(allowedOrigins) > 0 && allowedOrigins[0] == "*" {
				// If wildcard is configured, allow any origin
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
<<<<<<< HEAD
<<<<<<< HEAD
=======
		// // Add SuperTokens CORS headers
		// corsHeaders := supertokens.GetAllCORSHeaders()
		// for _, header := range corsHeaders {
		// 	c.Writer.Header().Set(header, "*") // Hoặc để trống nếu không cần value cụ thể
		// }
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
<<<<<<< HEAD
<<<<<<< HEAD
=======

func ContentSecurityPolicyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://cdn.jsdelivr.net; "+
				"style-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; "+
				"img-src 'self' data:; "+
				"font-src 'self' https://cdn.jsdelivr.net; "+
				"connect-src 'self' https://try.supertokens.com; "+
				"frame-ancestors 'self';",
		)
		c.Next()
	}
}
>>>>>>> ada3717 (Docker file)
=======
>>>>>>> 4ccd381 (Update appointment flow)
