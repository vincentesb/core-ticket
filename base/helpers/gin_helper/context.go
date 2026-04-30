package gin_helper

import (
	"core-ticket/base/helpers/base_helper"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Context interface {
	/*
		GinCtx returns the underlying gin.Context object associated with the ContextImpl instance.

		Returns:
		- *gin.Context: The gin.Context object associated with the ContextImpl instance.
	*/
	GinCtx() *gin.Context

	/*
		Identity returns the identity information extracted from the gin context.
		It retrieves the username, server code, database name, company code, company ID, and user role ID from the context.
		If any of the required information is missing or cannot be retrieved, it returns an empty Identity struct and an error.
		In case of an error during the retrieval process, it creates a new error with the UnknownError code.
		The method returns the extracted identity information and any error encountered during the process.

		Returns:
		- base_helper.Identity: The extracted identity information from the context.
		- error: An error indicating any issues encountered during the retrieval process.
	*/
	Identity() (base_helper.Identity, error)

	/*
		ShouldBind binds the passed object to the request data based on the request method and content type.
		It uses the Default method to determine the appropriate binding based on the method and content type of the request.

		Parameters:
		- obj: The object to bind the request data to.

		Returns:
		- error: An error if the binding process encounters any issues.
	*/
	ShouldBind(obj any) error

	/*
		ShouldBindQuery binds the passed object to the query parameters in the request.
		It uses the query binding to map the query parameters to the fields of the object.

		Parameters:
		- obj: The object to bind the query parameters to.

		Returns:
		- error: An error if the binding process fails.
	*/
	ShouldBindQuery(obj any) error

	/*
		ShouldBindForm binds the form data to the provided object using the form binding.
		It returns an error if the binding process fails.

		Parameters:
		- obj: The object to which the form data will be bound.

		Returns:
		- error: An error if the binding process fails.
	*/
	ShouldBindForm(obj any) error

	/*
		ShouldBindJSON binds the JSON data from the request body to the provided object.
		It uses the JSON binding method to perform the binding operation.

		Parameters:
		- obj: The object to which the JSON data will be bound.

		Returns:
		- error: An error if the binding operation fails.
	*/
	ShouldBindJSON(obj any) error

	/*
		ShouldBindUri binds the URI parameters from the gin context to the given object using the URI binding rules.

		Parameters:
		- obj: The object to bind the URI parameters to.

		Returns:
		- error: An error if the binding process fails.
	*/
	ShouldBindUri(obj any) error

	/*
		ShouldBindMultipartForm binds the multipart form data to the provided object using the form binding.
		It returns an error if the binding process fails.

		Parameters:
		- obj: The object to which the form data will be bound.

		Returns:
		- error: An error if the binding process fails.
	*/
	ShouldBindMultipartForm(obj any) error

	/*
		Default returns the appropriate binding based on the HTTP method and content type.

		Parameters:
		- method: a string representing the HTTP method (e.g., GET, POST)
		- contentType: a string representing the content type (e.g., application/json)

		Returns:
		- ginBinding.Binding: the binding object based on the method and contentType

		Example:

			context := &ContextImpl{}
			binding := context.Default("GET", "application/json")

		This method determines the binding to be used for parsing request data based on the HTTP method and content type.
	*/
	Default(method, contentType string) binding.Binding

	/*
		Validator returns the underlying *validator.Validate object associated with the ContextImpl instance.

		Returns:
		- *validator.Validate: The validator.Validate object associated with the ContextImpl instance.
	*/
	Validator() *validator.Validate
}
