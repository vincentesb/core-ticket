package middlewares

import (
	"core-ticket/base/helpers/http_helper"
	"core-ticket/base/token"
	"core-ticket/constants"
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (m *MiddlewareImpl) BasicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var basicToken = token.ExtractBasicToken(c)

		if basicToken == "" {
			http_helper.UnauthorizedResponse(c, constants.EC_UNAUTHORIZED, "Unauthorized", nil)
			return
		}

		decoded, errDecoded := base64.StdEncoding.DecodeString(basicToken)
		if errDecoded != nil {
			http_helper.UnauthorizedResponse(c, constants.EC_UNAUTHORIZED, "Unauthorized", nil)
			return
		}

		credentials := strings.SplitN(string(decoded), ":", 2)
		if len(credentials) != 2 {
			http_helper.UnauthorizedResponse(c, constants.EC_UNAUTHORIZED, "Unauthorized", nil)
			return
		}

		username := viper.GetString("BASIC_AUTH_USERNAME")
		password := viper.GetString("BASIC_AUTH_PASSWORD")

		if credentials[0] != username || credentials[1] != password {
			http_helper.UnauthorizedResponse(c, constants.EC_UNAUTHORIZED, "Unauthorized", nil)
			return
		}

		c.Next()
	}
}
