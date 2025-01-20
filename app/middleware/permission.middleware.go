package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/redis"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/quanganh247-qa/go-blog-be/app/util/perms"
)

func PermissionMiddleware(methods []perms.Permission) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the authenticated user's payload
		authPayload, err := GetAuthorizationPayload(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
			return
		}

		// Load user info from cache or database
		userInfo, err := redis.Client.UserInfoLoadCache(authPayload.Username)
		if err != nil {
			ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
			return
		}
		// log.Println("User info: ", userInfo.Role)

		isValid := perms.CheckPermission(methods, userInfo.Role)
		if !isValid {
			ctx.AbortWithStatusJSON(403, util.ErrorResponse(fmt.Errorf("Your account does not have permission to access the function [ %v ] ", methods)))
			return
		}

		// Proceed to the next handler
		ctx.Next()
	}
}
