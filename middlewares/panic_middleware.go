package middlewares

import (
	"core-ticket/base/helpers/http_helper"
	"fmt"

	"github.com/gin-gonic/gin"
)

func PanicMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic: %v\n", r)
				http_helper.HttpErrorResponse(c, fmt.Errorf("%v", r))
				return
			}
		}()

		c.Next()
	}
}
