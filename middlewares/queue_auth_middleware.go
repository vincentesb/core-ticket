package middlewares

import (
	"core-ticket/base/helpers/http_helper"
	"core-ticket/base/token"
	"core-ticket/constants"
	"core-ticket/constants/queue_engine"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (m *MiddlewareImpl) QueueAuthMiddleware(queueType queue_engine.QueueType) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestToken := token.ExtractJwtToken(c)
		tokenKey := queue_engine.EnvKeyReportApiToken
		if queueType == queue_engine.TypeUpload {
			tokenKey = queue_engine.EnvKeyUploadApiToken
		}

		queueApiToken := viper.GetString(tokenKey)
		if queueApiToken != requestToken {
			http_helper.UnauthorizedResponse(c, constants.EC_UNAUTHORIZED, "Unauthorized", nil)
			return
		}
		c.Next()
	}
}
