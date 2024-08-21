package middleware

// import (
// 	"fmt"

// 	"github.com/gin-gonic/gin"
// 	"github.com/hpt/go-client/app/service/redis"
// 	"github.com/hpt/go-client/app/util"
// 	"github.com/hpt/go-client/app/util/perms"
// )

// func PermissionMiddleware(methods []perms.Permission,typeApi perms.TypeApi) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authPayload, err := GetAuthorizationPayload(ctx)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
// 			return
// 		}
// 		userInfo, err := redis.Client.UserInfoLoadCache(authPayload.Username)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(403, util.ErrorResponse(err))
// 			return
// 		}

// 		isValid := perms.CheckPermission(methods, typeApi, &userInfo.Permissions)
// 		if !isValid {
// 			ctx.AbortWithStatusJSON(403, util.ErrorResponse(fmt.Errorf("Tài khoản của bạn không có quyền truy cập vào chức năng [ %v ] [ %s ]", methods,typeApi)))
// 			return
// 		}
// 		ctx.Next()
// 	}
// }
