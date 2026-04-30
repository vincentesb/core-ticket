package middlewares

import (
	"github.com/gin-gonic/gin"
)

func (m *MiddlewareImpl) TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
