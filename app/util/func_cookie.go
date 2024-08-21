package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetCookieSameSite(ctx *gin.Context) (host string, secure bool) {
	mode := Configs.CookieSameSite
	var sameSite http.SameSite
	switch mode {
	case "LAX":
		sameSite = http.SameSiteLaxMode
	case "STRICT":
		sameSite = http.SameSiteStrictMode
	case "NONE":
		sameSite = http.SameSiteNoneMode
	default:
		sameSite = http.SameSiteDefaultMode
	}
	ctx.SetSameSite(sameSite)
	host = ctx.Request.Host
	secure = Configs.CookieSecure
	if !Configs.CookieUseHost {
		host = ""
	}
	return
}
