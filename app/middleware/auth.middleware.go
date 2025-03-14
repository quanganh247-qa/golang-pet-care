package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/quanganh247-qa/go-blog-be/app/service/token"
	"github.com/quanganh247-qa/go-blog-be/app/util"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func GetAuthorizationPayload(ctx *gin.Context) (*token.Payload, error) {
	payload, exists := ctx.Get(AuthorizationPayloadKey)
	if !exists {
		return nil, errors.New("payload not found")
	}
	return payload.(*token.Payload), nil
}
<<<<<<< HEAD
=======

// AuthMiddleware kiểm tra phiên hợp lệ với SuperTokens
func SuperTokensAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Tạo biến bool và lấy địa chỉ của nó
		sessionRequired := true
		sessionContainer, err := session.GetSession(ctx.Request, ctx.Writer, &sessmodels.VerifySessionOptions{
			SessionRequired: &sessionRequired, // Sửa lỗi ở đây, // Bắt buộc phải có session
		})
		if err != nil {
			// Nếu không có session hoặc token không hợp lệ
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: " + err.Error(),
			})
			return
		}

		// Lưu thông tin session vào context để sử dụng trong handler
		ctx.Set("session", sessionContainer)
		ctx.Next()
	}
}
>>>>>>> ada3717 (Docker file)
