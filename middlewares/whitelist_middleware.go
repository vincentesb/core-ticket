package middlewares

import (
	"core-ticket/base/helpers/http_helper"

	"github.com/gin-gonic/gin"
)

func (m *MiddlewareImpl) WhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, err := m.IpWhitelistService.InWhitelist(c.ClientIP())
		if ok && err == nil {
			c.Next()
			return
		}

		http_helper.ForbiddenResponse(c, nil)
		return
	}
}
