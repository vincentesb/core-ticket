package gin_helper

import (
	"core-ticket/base/helpers/base_helper"
	"core-ticket/base/helpers/context_helper"
	"core-ticket/base/helpers/error_helper"
	cBinding "core-ticket/base/helpers/gin_helper/binding"
	"core-ticket/base/helpers/gin_helper/binding/form"
	"core-ticket/base/helpers/gin_helper/binding/json"
	"core-ticket/base/helpers/gin_helper/binding/query"
	"core-ticket/base/helpers/gin_helper/binding/uri"
	"core-ticket/base/helpers/gin_helper/binding/utility"
	"core-ticket/constants/error_code"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	ginBinding "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	formBinding  = form.Binding{}
	jsonBinding  = json.Binding{}
	queryBinding = query.Binding{}
	uriBinding   = uri.Binding{}
)

type ContextImpl struct {
	c *gin.Context
	v *validator.Validate
}

/*
ShouldBindMultipartForm binds the multipart form data to the provided object using the form binding.
It returns an error if the binding process fails.

Parameters:
- obj: The object to which the form data will be bound.

Returns:
- error: An error if the binding process fails.
*/
func (e *ContextImpl) ShouldBindMultipartForm(obj any) error {
	err := e.ShouldBind(obj)
	if err != nil {
		var numError *strconv.NumError
		switch {
		case errors.As(err, &numError),
			errors.Is(err, strconv.ErrRange),
			errors.Is(err, strconv.ErrSyntax):
			return utility.ValidateMap(
				utility.ConvertUrlValuesToMapStringInterface(e.c.Request.Form),
				utility.GetRulesFromStruct(obj, cBinding.Form),
			)
		}
		return err
	}

	return nil
}

/*
Validator returns the underlying *validator.Validate object associated with the ContextImpl instance.

Returns:
- *validator.Validate: The validator.Validate object associated with the ContextImpl instance.
*/
func (e *ContextImpl) Validator() *validator.Validate {
	return e.v
}

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
func (e *ContextImpl) Default(method, contentType string) ginBinding.Binding {
	if method == http.MethodGet {
		return formBinding
	}

	switch contentType {
	case ginBinding.MIMEJSON:
		return jsonBinding
	case ginBinding.MIMEMultipartPOSTForm:
		return ginBinding.FormMultipart
	default:
		return formBinding
	}
}

/*
ShouldBind binds the passed object to the request data based on the request method and content type.
It uses the Default method to determine the appropriate binding based on the method and content type of the request.

Parameters:
- obj: The object to bind the request data to.

Returns:
- error: An error if the binding process encounters any issues.
*/
func (e *ContextImpl) ShouldBind(obj any) error {
	return e.c.ShouldBindWith(obj, e.Default(e.c.Request.Method, e.c.ContentType()))
}

/*
ShouldBindQuery binds the passed object to the query parameters in the request.
It uses the query binding to map the query parameters to the fields of the object.

Parameters:
- obj: The object to bind the query parameters to.

Returns:
- error: An error if the binding process fails.
*/
func (e *ContextImpl) ShouldBindQuery(obj any) error {
	return e.c.ShouldBindWith(obj, queryBinding)
}

/*
ShouldBindForm binds the form data to the provided object using the form binding.
It returns an error if the binding process fails.

Parameters:
- obj: The object to which the form data will be bound.

Returns:
- error: An error if the binding process fails.
*/
func (e *ContextImpl) ShouldBindForm(obj any) error {
	return e.c.ShouldBindWith(obj, formBinding)
}

/*
ShouldBindJSON binds the JSON data from the request body to the provided object.
It uses the JSON binding method to perform the binding operation.

Parameters:
- obj: The object to which the JSON data will be bound.

Returns:
- error: An error if the binding operation fails.
*/
func (e *ContextImpl) ShouldBindJSON(obj any) error {
	return e.c.ShouldBindWith(obj, jsonBinding)
}

/*
ShouldBindUri binds the URI parameters from the gin context to the given object using the URI binding rules.

decodes an encoded URI parameters upon receiving request

Parameters:
- obj: The object to bind the URI parameters to.

Returns:
- error: An error if the binding process fails.
*/
func (e *ContextImpl) ShouldBindUri(obj any) error {
	m := make(map[string][]string)
	for _, v := range e.c.Params {
		decodedURI, err := url.QueryUnescape(v.Value)
		if err != nil {
			return err
		}
		m[v.Key] = []string{decodedURI}
	}
	return uriBinding.Bind(m, obj)
}

/*
Ctx returns a Context interface implementation with the provided gin.Context.

Parameters:
- c: The gin.Context to be used for the Context implementation.

Returns:
- Context: A Context interface implementation.

Example:

	ctx := Ctx(ginContext)
*/
func Ctx(c *gin.Context) Context {
	return &ContextImpl{c, validator.New()}
}

/*
GinCtx returns the underlying gin.Context object associated with the ContextImpl instance.

Returns:
- *gin.Context: The gin.Context object associated with the ContextImpl instance.
*/
func (e *ContextImpl) GinCtx() *gin.Context {
	return e.c
}

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
func (e *ContextImpl) Identity() (base_helper.Identity, error) {
	identity, err := context_helper.GetIdentity(e.c)
	if err != nil {
		return identity, error_helper.New(err, error_code.UnknownError)
	}

	return identity, nil
}
