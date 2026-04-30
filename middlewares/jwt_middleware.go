package middlewares

import (
	"github.com/gin-gonic/gin"
)

func (m *MiddlewareImpl) JwtMiddleware(isRefreshToken bool, companyAuth bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}
