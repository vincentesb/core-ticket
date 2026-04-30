package error_helper

import (
	"core-ticket/base/helpers/string_helper"
	"core-ticket/constants/error_code"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

/*
ValidationError represents an error that occurred during validation.
It contains information about the error code, tag, field, parameter, row, and message associated with the validation error.
*/
type ValidationError struct {
	ErrorCode error_code.ValidationErrorCode
	Tag       string
	Attribute string
	Field     string
	Param     string
	Row       int
	message   string
}

/*
Error returns the error message generated based on the validation error tag and field.

Returns:
  - string: The error message generated based on the validation error tag and field.
*/
func (ve *ValidationError) Error() string {
	return ve.generateErrorMessage()
}

/*
Set updates the ValidationError instance with the information from the provided validator.FieldError instance.
It sets the ErrorCode, Tag, Field, and Param fields of the ValidationError based on the corresponding values from the FieldError.
It then returns the updated ValidationError instance.

Parameters:
- fe: The validator.FieldError instance containing the error information to be set in the ValidationError.

Returns:
- *ValidationError: The updated ValidationError instance with the error information set.
*/
func (ve *ValidationError) Set(fe validator.FieldError) *ValidationError {
	return ve.config(fe)
}

/*
config updates the ValidationError instance with the error code, tag, field, and param from the provided validator.FieldError instance.

Parameters:
- fe: validator.FieldError - the field error instance containing information about the validation error

Returns:
- *ValidationError: the updated ValidationError instance
*/
func (ve *ValidationError) config(fe validator.FieldError) *ValidationError {
	ve.ErrorCode = ve.mapErrorCode(fe)
	ve.Tag = fe.Tag()
	ve.Field = fe.Field()
	ve.Param = fe.Param()
	ve.Attribute = fe.StructNamespace()
	ve.Row = ExtractRowFromNamespace(fe.Namespace())

	return ve
}

/*
mapErrorCode maps a validator FieldError to a corresponding ValidationErrorCode based on the validation tag.
If the tag is recognized, it returns the appropriate ValidationErrorCode. If the tag is not recognized, it returns ValidationError as the default error code.

Parameters:
- fe: validator FieldError to map to a ValidationErrorCode

Returns:
- ValidationErrorCode: the mapped ValidationErrorCode for the given FieldError
*/
func (ve *ValidationError) mapErrorCode(
	fe validator.FieldError,
) (errCode error_code.ValidationErrorCode) {
	switch fe.Tag() {
	case "required":
		errCode = error_code.IsRequired
	case "min":
		errCode = error_code.MinLength
	case "max":
		errCode = error_code.MaxLength
	case "email":
		errCode = error_code.Email
	case "phone_number":
		errCode = error_code.PhoneNumber
	case "url":
		errCode = error_code.Url
	case "integer":
		errCode = error_code.Integer
	case "float":
		errCode = error_code.Float
	case "boolean":
		errCode = error_code.IsBoolean
	case "date":
		errCode = error_code.Date
	case "in":
		errCode = error_code.In
	case "unique_default":
	case "unique_in_db":
		errCode = error_code.IsUnique
	case "exist_in_db":
		errCode = error_code.NotExistsInDB
	case "gt":
		errCode = error_code.GreaterThan
	case "gte":
		errCode = error_code.GreaterThanEqual
	case "lt":
		errCode = error_code.LessThan
	case "lte":
		errCode = error_code.LessThanEqual
	case "ne":
		errCode = error_code.NotEqual
	case "alpha":
		errCode = error_code.Alpha
	case "alphanum":
		errCode = error_code.AlphaNum
	case "alphanum_space":
		errCode = error_code.AlphaNumSpace
	case "alphanum_space_empty":
		errCode = error_code.AlphaNumSpace
	case "numeric":
		errCode = error_code.IsNumeric
	case "regex":
		errCode = error_code.Regex
	case "custom":
		errCode = error_code.Custom
	case "not_emoji":
		errCode = error_code.NotEmoji
	case "decimal":
		errCode = error_code.Decimal
	default:
		errCode = error_code.ValidationErrorCode(error_code.ValidationError)
	}

	return errCode
}

/*
generateErrorMessage generates an error message based on the validation error tag and field.

Returns:
- message (string): The error message generated based on the validation error.
*/
func (ve *ValidationError) generateErrorMessage() (message string) {
	switch ve.Tag {
	case "required":
		message = fmt.Sprintf("field %s is required and cannot be null", ve.Field)
	case "min":
		message = fmt.Sprintf("field %s cannot consist less than %v", ve.Field, ve.Param)
	case "max":
		message = fmt.Sprintf("field %s cannot consist more than %v", ve.Field, ve.Param)
	case "gt":
		message = fmt.Sprintf("%s value must greater than %s", ve.Field, ve.Param)
	case "gte":
		message = fmt.Sprintf("field %s must have value minimum of %v", ve.Field, ve.Param)
	case "lt":
		message = fmt.Sprintf("%s value must less than %s", ve.Field, ve.Param)
	case "lte":
		message = fmt.Sprintf("%s value must less than or equal %s", ve.Field, ve.Param)
	case "ne":
		message = fmt.Sprintf("%s value must not equal %s", ve.Field, ve.Param)
	case "exist_in_db":
		message = fmt.Sprintf("%s data is not found", ve.Field)
	case "unique_default":
		message = fmt.Sprintf("Data default has been taken by other data")
	case "date":
		message = fmt.Sprintf("%s does not implement %s date format", ve.Field, ve.Param)
	case "numeric":
		message = fmt.Sprintf("%s is not an numeric value", ve.Field)
	case "phone_number":
		message = fmt.Sprintf("%s value can only consist of +, -, and number", ve.Field)
	case "alpha":
		message = fmt.Sprintf("%s can only consist of alphabet", ve.Field)
	case "alphanum":
		message = fmt.Sprintf("%s can only consist of alphabet and number", ve.Field)
	case "alphanum_space", "alphanum_space_empty":
		message = fmt.Sprintf("%s can only consist of alphabet, number and space", ve.Field)
	case "unique_in_db":
		message = fmt.Sprintf("%s data already taken by other data", ve.Field)
	case "not_emoji":
		message = fmt.Sprintf("%s cannot contain emojis", ve.Field)
	case "npwp":
		message = fmt.Sprintf("%s must consist of 15 or 16 numeric digits", ve.Field)
	case "number":
		message = fmt.Sprintf("%s value can only consist of number", ve.Field)
	case "branch_location_valid":
		message = fmt.Sprintf("%s is not registered in branch", ve.Field)
	case "user_branch_valid":
		message = fmt.Sprintf("user doesn't have access to %s", ve.Field)
	case "customer_branch_valid":
		message = fmt.Sprintf("customer doesn't have access to %s", ve.Field)
	case "location_warehouse_kitchen_valid":
		message = fmt.Sprintf("%s type is not warehouse or kitchen", ve.Field)
	case "supplier_pic":
		message = fmt.Sprintf("there must be at least one default flag on the supplier PIC")
	case "unique":
		message = fmt.Sprintf("%s must be unique in %s", ve.Param, ve.Field)
	case "datetime":
		var format string
		switch ve.Param {
		case time.DateOnly:
			format = "YYYY-MM-DD"
		case time.RFC3339:
			format = "RFC3339/ISO8601: YYYY-MM-DDTHH:ii:ssZ"
		case "02/01/2006":
			format = "DD/MM/YYYY"
		default:
			format = ve.Param
		}
		message = fmt.Sprintf("%s does not implement '%s' date format", ve.Field, format)
	case "oneof":
		params := strings.Split(ve.Param, " ")
		formatted := params[0]

		if len(params) == 2 {
			formatted = strings.Join(params, " or ")
		} else if len(params) > 2 {
			formatted = fmt.Sprintf("%s or %s", strings.Join(params[:len(params)-1], ", "), params[len(params)-1])
		}

		message = fmt.Sprintf("%s value must be one of %s", ve.Field, formatted)
	case "required_unless":
		params := strings.Split(ve.Param, " ")
		message = fmt.Sprintf(
			"field %s is required and cannot be null unless %s is %s",
			ve.Field,
			params[0],
			params[1],
		)
	case "required_if":
		params := strings.Split(ve.Param, " ")
		message = fmt.Sprintf(
			"field %s is required and cannot be null when %s is %s",
			ve.Field,
			params[0],
			params[1],
		)
	case "excluded_unless":
		params := strings.Split(ve.Param, " ")
		message = fmt.Sprintf(
			"field %s is not required unless %s is %s",
			ve.Field,
			params[0],
			params[1],
		)
	case "excluded_if":
		params := strings.Split(ve.Param, " ")
		message = fmt.Sprintf(
			"field %s is not required when %s is %s",
			ve.Field,
			params[0],
			params[1],
		)
	case "not_idr_currency":
		message = fmt.Sprintf("field %s must be other than IDR currency", ve.Field)
	case "end_period":
		message = fmt.Sprintf("field %s must be greater than last close period date", ve.Field)
	case "in":
		message = fmt.Sprintf("field %s must be a valid data", ve.Field)
	case "status_valid":
		message = fmt.Sprintf("%s transaction must be %s", ve.Field, ve.Param)
	case "product_valid":
		message = fmt.Sprintf("%s is not valid on %s", ve.Field, ve.Param)
	case "user_access_valid":
		message = fmt.Sprintf("user doesn't have access to %s", ve.Field)
	case "custom":
		message = ve.message
	case "idr_currency_rate":
		message = fmt.Sprintf("field %s must be equal to 1", ve.Field)
	case "map_user_location":
		message = fmt.Sprintf("location is not registered with the user")
	case "decimal":
		message = fmt.Sprintf("%s must be shorter. Maximum is decimal(%s)", ve.Field, ve.Param)
	case "printascii":
		message = "only printable ASCII character are allowed."
	case "ascii":
		message = fmt.Sprintf(
			"only alphabet, number and symbols are allowed (`-=[]\\;\\',./~!@#$%%^&*()_+{}|:\"<>?) for %s",
			string_helper.ConvertCamelCaseToNormal(ve.Field),
		)
	case "required_if_cost_center":
		message = fmt.Sprintf("%s is required for users with cost center access", ve.Field)
	case "is_active":
		message = fmt.Sprintf("%s is not active, please reactivate to continue", ve.Field)
	case "compare_list":
		params := strings.Split(ve.Param, " ")
		src := params[0]
		dst := params[1]
		message = fmt.Sprintf("%s[%d].%s is not found in the %s", src, ve.Row, ve.Field, dst)
	case "request_template_product":
		message = "productDetailID is not valid on requestTemplateID"
	case "future_date":
		message = fmt.Sprintf("field %s must not be greater than today", ve.Field)
	case "product_detail_unique_per_bom":
		message = fmt.Sprintf("%s must be unique for the same BomID", ve.Field)

	default:
		message = fmt.Sprintf("failed to validate %s for the %s validation", ve.Field, ve.Tag)
	}
	return message
}

/*
NewValidationError creates a new instance of ValidationError with the provided parameters.

Parameters:
- errorCode: The error code for the validation error.
- tag: The tag associated with the validation error.
- field: The field name where the validation error occurred.
- param: The parameter related to the validation error.
- row: The row number where the validation error occurred.
- message: Optional message associated with the validation error.

Returns:
- A pointer to the newly created ValidationError instance.

Example:

	errCode := error_code.ValidationErrorCode(123)
	tag := "required"
	field := "username"
	param := "min_length"
	row := 1
	message := "Username must be at least 5 characters long"
	validationError := NewValidationError(errCode, tag, field, param, row, message)
*/
func NewValidationError(
	errorCode error_code.ValidationErrorCode,
	tag string,
	field string,
	param string,
	row int,
	message ...string,
) *ValidationError {
	ve := &ValidationError{
		ErrorCode: errorCode,
		Tag:       tag,
		Field:     field,
		Param:     param,
		Row:       row,
	}
	if message != nil {
		ve.message = message[0]
	}
	return ve
}

/*
ExtractRowFromNamespace to extract row detail index to give precise detail index num
*/
func ExtractRowFromNamespace(namespace string) int {
	var indexRegex = regexp.MustCompile(`\[(\d+)\]`)
	matches := indexRegex.FindAllStringSubmatch(namespace, -1)
	if len(matches) > 0 {
		// Use the last match (e.g., details[1])
		indexStr := matches[len(matches)-1][1]
		if idx, err := strconv.Atoi(indexStr); err == nil {
			return idx
		}
	}
	return 0
}
