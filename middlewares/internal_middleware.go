package middlewares

import (
	"core-ticket/base/helpers/http_helper"
	"core-ticket/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func InternalMiddleware(cfg config.AppConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]

		if privateHostname := cfg.AppPrivateHost; host != privateHostname {
			http_helper.ForbiddenResponse(c, nil)
			return
		}
		c.Next()
	}
}
