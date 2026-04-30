package middlewares

import (
	"core-ticket/base/helpers/http_helper"
	"core-ticket/constants"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TypeValidatorSource string

var (
	URIPathSource = TypeValidatorSource("uri")
	QuerySource   = TypeValidatorSource("query")
)

func TypeValidatorMiddleware[T any](dto T, source TypeValidatorSource) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			validationErrs = []http_helper.Error{}
			rt             = reflect.TypeOf(dto)
		)

		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			name := field.Tag.Get("typecheck")
			if name == "" {
				continue
			}

			var rawVal string
			switch source {
			case URIPathSource:
				rawVal = c.Param(name)
			case QuerySource:
				rawVal = c.Request.URL.Query().Get(name)
			}

			if rawVal == "" {
				continue
			}

			var err error
			switch field.Type.Kind() {
			case reflect.Bool:
				_, err = strconv.ParseBool(rawVal)
			case reflect.Int:
				_, err = strconv.Atoi(rawVal)
			case reflect.Int64:
				_, err = strconv.ParseInt(rawVal, 10, 64)
			case reflect.Float64:
				_, err = strconv.ParseFloat(rawVal, 64)
			default:
				continue
			}

			if err != nil {
				validationErrs = append(validationErrs, http_helper.Error{
					Code:      string(constants.EC_VALIDATION_ERROR),
					Message:   fmt.Sprintf("Invalid value for %s, value must be %s", name, field.Type.Kind().String()),
					Attribute: name,
				})
			}
		}

		if len(validationErrs) > 0 {
			http_helper.ErrorResponse(c, &http_helper.HandlerError{
				StatusCode: 400,
				RespCode:   string(constants.EC_VALIDATION_ERROR),
				Message:    "Validation error",
				Data:       validationErrs,
			})
			c.Abort()
			return
		}
	}

}
