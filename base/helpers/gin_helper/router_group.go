package gin_helper

import (
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/http_helper"
	"core-ticket/constants/error_code"
	"errors"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type HandlerFunc[T any] func(c Context) (T, error)

type RouterGroup struct {
	ginRouterGroup *gin.RouterGroup
}

/*
ginHandlerFunc is a function that takes a HandlerFunc as input and returns a gin.HandlerFunc.
It executes the provided handler function by passing a Context obtained from the gin.Context.
If an error occurs during execution, it generates an HTTP error response using http_helper.HttpErrorResponse.
If the execution is successful, it generates a success response using http_helper.SuccessResponse.

Parameters:
- handler: HandlerFunc[T] - The handler function to be executed.

Returns:
- gin.HandlerFunc: A function that can be registered as a route handler in Gin framework.
*/
func ginHandlerFunc[T any](handler HandlerFunc[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		response, err := handler(Ctx(c))
		if err != nil {
			var errs *error_helper.Error
			if errors.As(err, &errs) {
				http_helper.HttpErrorResponse(c, err, errs.Data())
			} else {
				http_helper.HttpErrorResponse(c, err)
			}
			return
		}

		var (
			rvResponse = reflect.ValueOf(response)
			message    = error_code.Success.Message()
			res        = any(response)
		)

		if res != nil && strings.Contains(rvResponse.Type().String(), "http_helper.HTTPResponse") {
			if !rvResponse.FieldByName("Message").IsZero() {
				message = rvResponse.FieldByName("Message").String()
			}
			res = rvResponse.FieldByName("Result").Interface()
		}

		http_helper.SuccessResponse(c, error_code.Success.String(), message, res)
	}
}

/*
GET registers a GET request handler for the specified relative path with the given RouterGroup.

Parameters:
- rg: Pointer to the RouterGroup where the GET request handler will be registered.
- relativePath: The relative path at which the GET request handler will be registered.
- handler: The handler function that processes the GET request and returns a response of type T.

Example:

	GET(&routerGroup, "/users", func(c Context) (User, error) {
		// Handler logic here
	})

Note:

	The handler function must conform to the HandlerFunc signature: func(c Context) (T, error).
*/
func GET[T any](rg *RouterGroup, relativePath string, handler HandlerFunc[T]) {
	rg.ginRouterGroup.GET(relativePath, ginHandlerFunc(handler))
}

/*
POST registers a POST request handler for the given relative path with the provided handler function.

Parameters:
- rg: Pointer to a RouterGroup struct representing the router group where the handler will be registered.
- relativePath: String representing the relative path at which the handler will be registered.
- handler: HandlerFunc[T] function that defines the logic for handling the POST request.

Example:

	POST(&RouterGroup{}, "/example", func(c Context) (T, error) {
		// Handler logic here
	})

Note:
- The handler function should have the signature func(c Context) (T, error), where T is the expected response type.
*/
func POST[T any](rg *RouterGroup, relativePath string, handler HandlerFunc[T]) {
	rg.ginRouterGroup.POST(relativePath, ginHandlerFunc(handler))
}

/*
PUT registers a PUT request handler for the specified relative path with the given handler function.

Parameters:
- rg: Pointer to a RouterGroup struct where the PUT request handler will be registered.
- relativePath: The relative path at which the PUT request handler will be registered.
- handler: HandlerFunc[T] function that defines the logic for handling the PUT request.

Example:

	PUT(&RouterGroup, "/example", func(c Context) (T, error) {
		// PUT request handling logic here
	})

Note:
- The handler function should have the signature func(c Context) (T, error), where T is the type of response data and error is the potential error returned.
*/
func PUT[T any](rg *RouterGroup, relativePath string, handler HandlerFunc[T]) {
	rg.ginRouterGroup.PUT(relativePath, ginHandlerFunc(handler))
}

/*
PATCH registers a PATCH request handler for the specified relative path with the given handler function.

Parameters:
- rg: Pointer to a RouterGroup struct representing the router group where the PATCH request handler will be registered.
- relativePath: String representing the relative path at which the PATCH request handler will be registered.
- handler: HandlerFunc[T] function that defines the logic to be executed when the PATCH request is received.

Example:

	PATCH(&RouterGroup, "/path", func(c Context) (T, error) {
		// handler logic here
	})

Note:

	The handler function should conform to the HandlerFunc[T] type signature.
*/
func PATCH[T any](rg *RouterGroup, relativePath string, handler HandlerFunc[T]) {
	rg.ginRouterGroup.PATCH(relativePath, ginHandlerFunc(handler))
}

/*
DELETE registers a DELETE request handler for the given relative path with the provided handler function.

Parameters:
- rg: Pointer to a RouterGroup struct which contains the gin RouterGroup.
- relativePath: The relative path at which the DELETE request will be registered.
- handler: HandlerFunc function that defines the logic to be executed when the DELETE request is received.

Returns:
- None
*/
func DELETE[T any](rg *RouterGroup, relativePath string, handler HandlerFunc[T]) {
	rg.ginRouterGroup.DELETE(relativePath, ginHandlerFunc(handler))
}

/*
GinRouterGroup returns the underlying gin.RouterGroup instance associated with the RouterGroup.
This allows access to the full range of methods and properties provided by the gin framework for routing and middleware handling.

Returns:
- *gin.RouterGroup: The gin.RouterGroup instance associated with the RouterGroup.
*/
func (rg *RouterGroup) GinRouterGroup() *gin.RouterGroup {
	return rg.ginRouterGroup
}
