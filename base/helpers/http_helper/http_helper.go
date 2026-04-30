package http_helper

import (
	"core-ticket/base/helpers/array_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/sentry_helper"
	"core-ticket/constants"
	"core-ticket/constants/error_code"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/iancoleman/strcase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Path      string      `json:"path"`
	TimeStamp string      `json:"timestamp"`
	Status    string      `json:"status"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
	Errors    []Error     `json:"errors"`
}

type Error struct {
	Attribute string `json:"attribute"`
	Field     string `json:"field,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Row       int    `json:"row,omitempty"`
}

type HandlerError struct {
	StatusCode int
	RespCode   string
	Message    string
	Data       []Error
}

type PaginationResult[T interface{}] struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Count int    `json:"count"`
	Data  []T    `json:"data"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}

type HTTPResponse[T any] struct {
	Result  T
	Message string
}

/*
NewHttpResponse creates a new HTTP response object with the provided parameters.

Parameters:
- request: *http.Request - The HTTP request object.
- statusCode: int - The status code of the response.
- respCode: string - The response code.
- message: string - The message associated with the response.
- result: interface{} - The result data of the response.
- errors: []Error - A slice of Error objects representing any errors in the response.

Returns:
- *Response: A pointer to the Response object created.

Deprecated: This function is deprecated and should not be used anymore.
*/
func NewHttpResponse(request *http.Request, statusCode int, respCode string, message string, result interface{}, errors []Error) *Response {
	if statusCode == 0 {
		statusCode = 200
	}
	if respCode == "" {
		respCode = constants.EC_SUCCESS
	}
	path := fmt.Sprintf("%s://%s%s", getScheme(request), request.Host, request.RequestURI)
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	status := "fail"
	if isHttpOk(statusCode) {
		status = "ok"
	}

	return &Response{
		Path:      path,
		TimeStamp: formattedTime,
		Status:    status,
		Code:      respCode,
		Message:   message,
		Result:    result,
		Errors:    errors,
	}
}

/*
NewHandlerError creates and returns a new instance of HandlerError with the provided parameters.

Parameters:
- statusCode: an integer representing the status code of the error.
- respCode: a string representing the response code of the error.
- message: a string containing the error message.
- data: a slice of Error structs containing additional error details.

Returns:
- a pointer to the newly created HandlerError instance.

Deprecated: This function is deprecated and should not be used anymore.
*/
func NewHandlerError(statusCode int, respCode string, message string, data []Error) *HandlerError {
	return &HandlerError{
		StatusCode: statusCode,
		RespCode:   respCode,
		Message:    message,
		Data:       data,
	}
}

/*
SuccessResponse generates a successful HTTP response with the provided response code, message, and result. If the message is empty, it defaults to "Ok".
It creates a new HTTP response using the NewHttpResponse function and sends it as a JSON response with a status code of 200.
The function then aborts the current request handling and returns.

Parameters:
- c (*gin.Context): The Gin context object for the current HTTP request.
- respCode (string): The response code to be included in the response.
- message (string): The message to be included in the response.
- result (interface{}): The result data to be included in the response.

Returns:
- None
*/
func SuccessResponse(c *gin.Context, respCode string, message string, result interface{}) {
	if message == "" {
		message = "Ok"
	}
	res := NewHttpResponse(c.Request, 200, respCode, message, result, nil)
	c.JSON(200, res)
	c.Abort()
	return
}

/*
PaginationResponse generates a PaginationResult object based on the provided parameters.

Parameters:
- c: *gin.Context - The Gin context object for the current request.
- result: []T - The slice of generic type data to be paginated.
- count: int - The total count of items in the dataset.
- page: int - The current page number.
- limit: int - The limit of items per page.

Returns:
- PaginationResult[T] - The paginated result containing page information, data slice, previous and next page URLs.

Example:
paginationResponse := PaginationResponse(c, dataSlice, totalCount, currentPage, itemsPerPage)
*/
func PaginationResponse[T any](c *gin.Context, result []T, count, page, limit int) PaginationResult[T] {
	paginationResult := PaginationResult[T]{}
	path := fmt.Sprintf("%s://%s%s", getScheme(c.Request), c.Request.Host, c.Request.RequestURI)

	nextUrl := ""
	if len(result) >= limit {
		if (page*limit)-count != 0 {
			nextPage := page + 1
			nextUrl = createPaginationPath(path, nextPage, limit)
		}
	}

	prevUrl := ""
	if page > 1 {
		prevPage := page - 1
		prevUrl = createPaginationPath(path, prevPage, limit)
	}

	paginationResult = PaginationResult[T]{
		Page:  page,
		Limit: limit,
		Count: count,
		Data:  result,
		Prev:  prevUrl,
		Next:  nextUrl,
	}

	return paginationResult
}

/*
SuccessPaginationResponse generates a success response with pagination information.

Parameters:
- c: gin.Context - The gin context for the request
- respCode: string - The response code for the success response
- message: string - The message to be included in the response
- result: []T - The data to be included in the response
- count: int - The total count of items
- page: int - The current page number
- limit: int - The limit of items per page

Returns:
- None

Example:
SuccessPaginationResponse(c, "200", "Success", data, totalCount, currentPage, itemsPerPage)
*/
func SuccessPaginationResponse[T interface{}](c *gin.Context, respCode string, message string, result []T,
	count int, page int, limit int) {
	SuccessResponse(c, respCode, message, PaginationResponse[T](c, result, count, page, limit))
	return
}

/*
createPaginationPath generates a pagination path based on the provided path, page number, and limit.
It concatenates the path with query parameters for page and limit using the provided values.

Parameters:
- path (string): The base path to which the pagination query parameters will be added.
- page (int): The page number to be included in the pagination path.
- limit (int): The limit of items per page to be included in the pagination path.

Returns:
- string: The pagination path with the added query parameters for page and limit.

Example:
paginationPath := createPaginationPath("/api/items", 2, 10)
// Result: "/api/items?page=2&limit=10"
*/
func createPaginationPath(path string, page int, limit int) string {
	return path + "?page=" + strconv.Itoa(page) + "&limit=" + strconv.Itoa(limit)
}

/*
ErrorResponse generates an HTTP response with the provided error details.

Parameters:
- c: The gin context to generate the response.
- err: The HandlerError containing the error details.

Deprecated: This function is deprecated and should not be used anymore.
*/
func ErrorResponse(c *gin.Context, err *HandlerError) {
	res := NewHttpResponse(c.Request, err.StatusCode, err.RespCode, err.Message, nil, err.Data)
	c.JSON(err.StatusCode, res)
	c.Abort()
	return
}

/*
HttpErrorResponse generates an HTTP response based on the provided error and result data.

Parameters:
- c: *gin.Context - The Gin context for the HTTP request.
- err: error - The error that occurred.
- result: ...interface{} - Optional result data to include in the response.

Returns:
- None

Behavior:
- If the error is of type *error_helper.Error, it constructs a response based on the error details and any provided result data.
- If the error is a validation error, it includes details of each validation error in the response.
- If the error is not of type *error_helper.Error, it generates a generic unknown error response.
- The function also captures exceptions using Sentry if available in the context.
- Finally, it sends the constructed response as JSON and aborts further processing of the request.
*/
func HttpErrorResponse(c *gin.Context, err error, result ...interface{}) {
	var (
		errs    *error_helper.Error
		res     *Response
		errCode int
	)

	if errors.As(err, &errs) {
		if errs.ErrorCode() == error_code.ValidationError {
			var he []Error
			for _, ve := range errs.ValidationErrors() {
				he = append(he, Error{
					Attribute: ve.Field,
					Field:     ve.Attribute,
					Code:      ve.ErrorCode.String(),
					Message:   ve.Error(),
					Row:       ve.Row,
				})
			}
			if result != nil {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), result[0], he)
			} else {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), nil, he)
			}
			errCode = errs.ErrorCode().HttpStatusCode()
		} else if array_helper.InArray([]error_code.ErrorCode{error_code.NotFound, error_code.Forbidden}, errs.ErrorCode()) {
			if result != nil {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), result[0], nil)
			} else {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), nil, nil)
			}
			errCode = errs.ErrorCode().HttpStatusCode()
		} else {
			if result != nil {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), result[0], nil)
			} else {
				res = NewHttpResponse(c.Request, errs.ErrorCode().HttpStatusCode(), errs.ErrorCode().String(), errs.Error(), nil, nil)
			}
			errCode = errs.ErrorCode().HttpStatusCode()

			sentry_helper.CaptureException(c, err)
		}
	} else {
		if result != nil {
			res = NewHttpResponse(c.Request, http.StatusInternalServerError, error_code.UnknownError.String(), "Unknown Error", result[0], nil)
		} else {
			res = NewHttpResponse(c.Request, http.StatusInternalServerError, error_code.UnknownError.String(), "Unknown Error", nil, nil)
		}
		errCode = http.StatusInternalServerError

		sentry_helper.CaptureException(c, err)
	}
	c.JSON(errCode, res)
	c.Abort()
	return
}

/*
BadRequestResponse generates a bad request response with the provided response code, message, errors, and optional result. If the message is empty, it defaults to "Bad Request". The function constructs a Response object using the NewHttpResponse function and sends it as a JSON response with a status code of http.StatusBadRequest. Finally, it aborts the current request.

Parameters:
- c (*gin.Context): The Gin context object.
- respCode (string): The response code for the bad request.
- message (string): The message to be included in the response.
- errors ([]Error): A slice of Error objects representing any errors.
- result (...interface{}): Optional result data to be included in the response.

Deprecated: This function is deprecated and should not be used anymore.

Returns:
- None
*/
func BadRequestResponse(c *gin.Context, respCode string, message string, errors []Error, result ...interface{}) {
	if message == "" {
		message = "Bad Request"
	}
	var res *Response
	if result != nil {
		res = NewHttpResponse(c.Request, 400, respCode, message, result[0], errors)
	} else {
		res = NewHttpResponse(c.Request, 400, respCode, message, nil, errors)
	}
	c.JSON(http.StatusBadRequest, res)
	c.Abort()
	return
}

/*
ServerErrorResponse generates a server error response with the given response code, message, and errors. If the message is empty, it defaults to "Internal Server Error". It creates an HTTP response using the NewHttpResponse function and sends it as a JSON response with a status code of 500 (Internal Server Error).

Parameters:
- c: *gin.Context - The Gin context for the request.
- respCode: string - The response code to be included in the response.
- message: string - The message to be included in the response.
- respErrors: []Error - A slice of Error structs representing any errors to be included in the response.

Returns:
- None

Deprecated: This function is deprecated and should not be used anymore.

Example:

	ServerErrorResponse(c, "SOME_ERROR_CODE", "An error occurred", []Error{Error{Attribute: "attr1", Code: "ERR_CODE_1", Message: "Error 1"}})
*/
func ServerErrorResponse(c *gin.Context, respCode string, message string, respErrors []Error) {
	if message == "" {
		message = "Internal Server Error"
	}
	res := NewHttpResponse(c.Request, 500, respCode, message, nil, respErrors)

	for _, err := range respErrors {
		sentry_helper.CaptureException(c, errors.New(err.Message))
	}

	c.JSON(http.StatusInternalServerError, res)
	c.Abort()
	return
}

/*
UnauthorizedResponse generates an HTTP response with status code 401 (Unauthorized) containing the provided response code, message, and errors. If the message is empty, a default "Unauthorized" message is used. The function creates a new HttpResponse object and returns it as a JSON response. It then aborts the current request handling.

Parameters:
- c (*gin.Context): The Gin context object for the current HTTP request.
- respCode (string): The response code to be included in the response.
- message (string): The message to be included in the response.
- errors ([]Error): A slice of Error objects to be included in the response.

Returns:
- None

Deprecated: This function is deprecated and should not be used anymore.

Example:
UnauthorizedResponse(c, "UNAUTHORIZED_ACCESS", "User is not authorized to access this resource", []Error{Error{Attribute: "permission", Code: "403", Message: "Insufficient permissions"}})
*/
func UnauthorizedResponse(c *gin.Context, respCode string, message string, errors []Error) {
	if message == "" {
		message = "Unauthorized"
	}
	res := NewHttpResponse(c.Request, 401, respCode, message, nil, errors)
	c.JSON(http.StatusUnauthorized, res)
	c.Abort()
	return
}

/*
ForbiddenResponse generates a response with a 403 status code (Forbidden) containing the provided errors.
It constructs the response using the NewHttpResponse function and sends it back to the client.
The function then aborts the request processing.

Parameters:
- c: The gin context object for the current HTTP request
- errors: A slice of Error structs containing details about the errors encountered

Returns: None

Deprecated: This function is deprecated and should not be used anymore.
*/
func ForbiddenResponse(c *gin.Context, errors []Error) {
	message := "Forbidden"
	respCode := constants.EC_FORBIDDEN
	res := NewHttpResponse(c.Request, 403, respCode, message, nil, errors)
	c.JSON(http.StatusForbidden, res)
	c.Abort()
	return
}

/*
NotFoundResponse generates a response for a "Not Found" error with the provided message and errors.

Parameters:
- c: The gin context object for the request
- message: The message to be included in the response. If empty, defaults to "Not Found"
- errors: A slice of Error structs containing details about any errors encountered

Returns:
- None

Deprecated: This function is deprecated and should not be used anymore.

Example:
NotFoundResponse(c, "Resource not found", []Error{Error{Attribute: "ID", Code: "EC03100004", Message: "Resource ID not found"}})
*/
func NotFoundResponse(c *gin.Context, message string, errors []Error) {
	if message == "" {
		message = "Not Found"
	}
	respCode := constants.EC_NOT_FOUND
	res := NewHttpResponse(c.Request, 404, respCode, message, nil, errors)
	c.JSON(http.StatusNotFound, res)
	c.Abort()
	return
}

/*
getScheme returns the scheme (http or https) based on the provided HTTP request object.
It checks if the request is using TLS (Transport Layer Security) and returns "https" in that case.
If the request header "X-Forwarded-Proto" is set to "https", it also returns "https".
Otherwise, it defaults to "http".

Parameters:
- r (*http.Request): The HTTP request object for which the scheme needs to be determined.

Returns:
- string: The scheme ("http" or "https") based on the request's TLS status and header.

Example:
scheme := getScheme(request)
// Result: "https" or "http"
*/
func getScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}

	if r.Header.Get("X-Forwarded-Proto") == "https" {
		return "https"
	}

	return "http"
}

/*
isHttpOk checks if the provided HTTP status code indicates a successful response.
It returns true if the status code is in the range of 200 (inclusive) to 299 (exclusive), indicating a successful HTTP response.
Otherwise, it returns false.

Parameters:
- statusCode (int): The HTTP status code to be checked.

Returns:
- bool: true if the status code represents a successful response, false otherwise.

Example:
isSuccess := isHttpOk(200)
// Result: true
*/
func isHttpOk(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < 300
}

/*
ParseBindErrors takes an error as input and returns a list of Error structs. It inspects the type of error and creates an Error struct accordingly. If the error is of type mysql.MySQLError, it creates a ValidationError with the error details. If the error is of type validator.ValidationErrors, it iterates over each validation error and creates a ValidationError for each. For any other type of error, it creates a generic ValidationError with the error message. The function then converts the ValidationError to an Error struct and appends it to the list of errors to be returned.

Parameters:
- err: The error to be parsed and converted into a list of Error structs.

Returns:
- []Error: A list of Error structs representing the parsed errors.

Deprecated: This function is deprecated and should not be used anymore.
*/
func ParseBindErrors(err error) []Error {
	var errs []Error
	var mySQLError *mysql.MySQLError
	var validationErrors validator.ValidationErrors
	var ve *error_helper.ValidationError
	switch {
	case errors.As(err, &mySQLError):
		ve = error_helper.NewValidationError(
			error_code.ValidationErrorCode(error_code.UnknownError),
			"mysql",
			mySQLError.Message,
			"",
			0,
		)
	case errors.As(err, &validationErrors):
		for _, fe := range err.(validator.ValidationErrors) {
			ve = error_helper.NewValidationError(
				error_code.ValidationErrorCode(error_code.ValidationError),
				fe.Tag(),
				fe.Field(),
				fe.Param(),
				0,
			)
		}
	default:
		ve = error_helper.NewValidationError(
			error_code.ValidationErrorCode(error_code.UnknownError),
			"unknown",
			err.Error(),
			"",
			0,
		)
	}
	validationError := Error{
		Attribute: ve.Field,
		Code:      ve.ErrorCode.String(),
		Message:   ve.Error(),
		Row:       ve.Row,
	}
	errs = append(errs, validationError)
	return errs
}

/*
ValidationErrorResponse generates a bad request response with validation errors or unknown errors.

Parameters:
- c: gin.Context: The gin context for the HTTP request
- err: error: The error to be checked for validation errors

Returns:
- None

Deprecated: This function is deprecated and should not be used anymore.

Example:

	err := validateInput(input)
	if err != nil {
		ValidationErrorResponse(c, err)
	}
*/
func ValidationErrorResponse(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		httpErrors := make([]Error, 0, len(validationErrors))
		for _, ve := range validationErrors {
			httpErrors = append(httpErrors, Error{
				Attribute: ve.Field(),
				Code:      constants.EC_VALIDATION_ERROR,
				Message:   ve.Error() + " with value:" + fmt.Sprintf("%v", ve.Value()),
			})
		}
		BadRequestResponse(c, constants.EC_VALIDATION_ERROR, "Error Validation", httpErrors)
	} else {
		BadRequestResponse(c, constants.EC_UNKNOWN_ERROR, err.Error(), nil)
	}
}

/*
MappingErrorMessage generates an error message based on the validation error provided.
It inspects the validation error tag and field to construct a human-readable error message.
The function handles various validation error cases such as required fields, maximum length, numeric values, uniqueness, etc.

Parameters:
- e (validator.FieldError): The validation error object containing information about the field validation failure.

Returns:
- string: The human-readable error message corresponding to the validation error.

Deprecated: This function is deprecated and should not be used anymore.

Example:
errorMessage := MappingErrorMessage(e)
// Result: "field <field> is required and cannot be null"

Note: This function is used to provide detailed error messages for validation failures.
*/
func MappingErrorMessage(e validator.FieldError) (message string) {
	switch e.Tag() {
	case "required":
		message = fmt.Sprintf("field %s is required and cannot be null", e.Field())
	case "max":
		message = fmt.Sprintf("field %s cannot consist more than %v", e.Field(), e.Param())
	case "gte":
		message = fmt.Sprintf("field %s must have value minimum of %v", e.Field(), e.Param())
	case "exists_in_db":
		message = fmt.Sprintf("%s is not found", e.Field())
	case "unique_default":
		message = fmt.Sprintf("%s does not implement d-m-Y date format", e.Value())
	case "numeric":
		message = fmt.Sprintf("%s is not an numeric value", e.Value())
	case "gt":
		message = fmt.Sprintf("%s value must greater than %s", e.Field(), e.Param())
	case "alphanum_space":
		message = fmt.Sprintf("%s can only consist of alphabet,number and space", e.Field())
	case "unique":
		message = fmt.Sprintf("%s with value %s is already taken", e.Field(), e.Value())
	case "required_with":
		message = fmt.Sprintf("%s is required with field %s", strcase.ToLowerCamel(e.Field()), strcase.ToLowerCamel(e.Param()))
	case "before_end_period":
		message = fmt.Sprintf("%s must be greater than end period date", strcase.ToLowerCamel(e.Field()))
	default:
		message = fmt.Sprintf("failed to validate %s for the %s validation", e.Field(), e.Tag())
	}
	return message
}
