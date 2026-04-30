package struct_validator

import (
	"bytes"
	"core-ticket/base/helpers/array_helper"
	"core-ticket/base/helpers/error_helper"
	"core-ticket/base/helpers/http_helper"
	"core-ticket/base/utility/nullable"
	"core-ticket/constants"
	"core-ticket/constants/error_code"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	"gopkg.in/guregu/null.v4"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	gout "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

/*
Deprecated: this var is deprecated and should not be used.
*/
var validate *Validate

/*
Deprecated: this struct is deprecated and should not be used.
*/
type Validate struct {
	validatorValidate *validator.Validate
	trans             gout.Translator
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func New() *Validate {
	if validate != nil {
		return validate
	}

	var trans gout.Translator
	validatorValidate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// To prevent error undefined trim validation
		_ = validatorValidate.RegisterValidation(
			"trim",
			func(fl validator.FieldLevel) bool { return true },
		)
		_ = validatorValidate.RegisterValidation("alphanum_space", ValidateAlphaNumericSpaceOnly)
		_ = validatorValidate.RegisterValidation(
			"alphanum_space_empty",
			ValidateAlphaNumericSpaceOrEmpty,
		)
		_ = validatorValidate.RegisterValidation("date", ValidateDate)
		_ = validatorValidate.RegisterValidation("not_emoji", NotEmoji())
		_ = validatorValidate.RegisterValidation("in", ValidateInArray)
		_ = validatorValidate.RegisterValidation("decimal", ValidateDecimal)

		// Replace field name using the json, form, or label tag registered in the struct
		validatorValidate.RegisterTagNameFunc(func(structField reflect.StructField) string {
			name := validate.getStructName(structField)
			// skip if tag key says it should be ignored
			if name == "-" {
				return ""
			}
			return name
		})

		trans, _ = gout.New(en.New(), en.New()).GetTranslator("en")
	}

	validate = &Validate{
		validatorValidate: validatorValidate,
		trans:             trans,
	}

	return validate
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) Struct(c *gin.Context, s interface{}) *http_helper.HandlerError {
	if err := v.validatorValidate.Struct(s); err != nil {
		return v.validationErrorToHttpError(err, BindingJSON, s)
	}

	v.trim(s)

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) ParseJSON(c *gin.Context, s interface{}) *http_helper.HandlerError {
	if err := v.validateContentType(c.Request, "application/json"); err != nil {
		return err
	}

	if httpErr := v.validateRuleWith(c, BindingJSON, s); httpErr != nil {
		return httpErr
	}

	v.trim(s)

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) ParseQuery(c *gin.Context, s interface{}) *http_helper.HandlerError {
	if httpErr := v.validateRuleWith(c, BindingQuery, s); httpErr != nil {
		return httpErr
	}

	v.trim(s)

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) ParseUri(c *gin.Context, s interface{}) *http_helper.HandlerError {
	if httpErr := v.validateRuleWith(c, BindingUri, s); httpErr != nil {
		return httpErr
	}

	v.trim(s)

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) validateRuleWith(
	c *gin.Context,
	b Binding,
	s interface{},
) *http_helper.HandlerError {
	if msg := v.validateNumericFromRequest(c, b, s); len(msg) > 0 {
		httpErr := http_helper.NewHandlerError(
			http.StatusBadRequest,
			constants.EC_IS_NUMERIC,
			"Parameter must be numeric",
			nil,
		)
		return httpErr
	}

	switch b {
	case BindingJSON, BindingQuery:
		binding := v.stringToHttpBinding(b)
		if err := c.ShouldBindWith(s, binding); err != nil {
			return v.validationErrorToHttpError(err, b, s)
		}
	case BindingUri:
		if err := c.ShouldBindUri(s); err != nil {
			return v.validationErrorToHttpError(err, b, s)
		}
	default:
		httpErr := http_helper.NewHandlerError(
			http.StatusBadRequest,
			constants.EC_VALIDATION_ERROR,
			"binding validation not implemented",
			nil,
		)
		return httpErr
	}

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) validationErrorToHttpError(
	err error,
	b Binding,
	s interface{},
) *http_helper.HandlerError {
	var errs []http_helper.Error
	for _, fe := range err.(validator.ValidationErrors) {
		ve := error_helper.NewValidationError(
			error_code.ValidationErrorCode(error_code.ValidationError),
			fe.Tag(),
			fe.Field(),
			fe.Param(),
			0,
		)

		validationError := http_helper.Error{
			Attribute: ve.Field,
			Code:      ve.ErrorCode.String(),
			Message:   ve.Error(),
		}
		errs = append(errs, validationError)
	}

	httpErr := http_helper.NewHandlerError(
		http.StatusBadRequest,
		constants.EC_VALIDATION_ERROR,
		"Validation Error HTTP",
		errs,
	)
	return httpErr
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) validateContentType(
	r *http.Request,
	contentType string,
) *http_helper.HandlerError {
	if ct := r.Header.Get("Content-Type"); ct != contentType {
		for _, v := range strings.Split(ct, ",") {
			mt, _, err := mime.ParseMediaType(v)
			if err != nil {
				httpErr := http_helper.NewHandlerError(
					http.StatusBadRequest,
					constants.EC_VALIDATION_ERROR,
					"invalid content type",
					nil,
				)
				return httpErr
			}
			if mt != contentType {
				httpErr := http_helper.NewHandlerError(
					http.StatusBadRequest,
					constants.EC_VALIDATION_ERROR,
					fmt.Sprintf("expecting %s, but got %s", contentType, ct),
					nil,
				)
				return httpErr
			}
		}
	}

	return nil
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) validateNumericFromRequest(
	c *gin.Context,
	b Binding,
	s interface{},
) map[string]string {
	body := v.requestValueToMap(c, b)
	val := reflect.ValueOf(s).Elem()

	return v.validateNumericInStruct(val, b, body, make(map[string]string))
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) requestValueToMap(c *gin.Context, b Binding) map[string]interface{} {
	body := make(map[string]interface{})

	switch b {
	case BindingQuery:
		rawQuery := c.Request.URL.Query()
		for k, v := range rawQuery {
			body[k] = v[0]
		}
	case BindingUri:
		rawParams := c.Params
		for _, param := range rawParams {
			body[param.Key] = param.Value
		}
	default:
		rawJSON, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(rawJSON))
		_ = json.Unmarshal(rawJSON, &body)
	}

	return body
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) validateNumericInStruct(
	rf reflect.Value,
	b Binding,
	body map[string]interface{},
	msg map[string]string,
) map[string]string {
	for i := 0; i < rf.NumField(); i++ {
		field := rf.Field(i)
		if field.Kind() == reflect.Struct {
			msg = v.validateNumericInStruct(rf.Field(i), b, body, msg)
		}

		structField := rf.Type().Field(i)
		key := v.getStructKey(structField, b)
		tag := structField.Tag.Get("binding")
		if strings.Contains(tag, "numeric") || strings.Contains(tag, "number") {
			keyVal, ok := body[key]
			if !ok || keyVal == "" {
				break
			}
			if _, err := strconv.Atoi(fmt.Sprint(body[key])); err != nil {
				msg[key] = fmt.Sprintf("%s harus berupa angka", key)
			}
		}
	}

	return msg
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) GetValidator() *validator.Validate {
	return v.validatorValidate
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) trim(s interface{}) {
	val := reflect.ValueOf(s).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := val.Type().Field(i)
		tag := structField.Tag.Get("binding")
		if strings.Contains(tag, "trim") {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) getStructName(structField reflect.StructField) string {
	name := strings.ToLower(structField.Name)
	if validate.getStructKey(structField, "form") != "" {
		name = validate.getStructKey(structField, "form")
	}
	if validate.getStructKey(structField, "form") != "" {
		name = validate.getStructKey(structField, "form")
	}
	if validate.getStructKey(structField, "json") != "" {
		name = validate.getStructKey(structField, "json")
	}
	if validate.getStructKey(structField, "label") != "" {
		name = validate.getStructKey(structField, "label")
	}

	return name
}

/*
Deprecated: this function is deprecated and should not be used.
*/
func (v *Validate) getStructKey(structField reflect.StructField, b Binding) string {
	return strings.SplitN(structField.Tag.Get(string(b)), ",", 2)[0]
}

func NullFlake(field reflect.Value) interface{} {
	switch field.Interface().(type) {
	case null.String:
		if nullableString, ok := field.Interface().(null.String); ok {
			if nullableString.Valid {
				return nullableString.String
			}
			return ""
		}
	case null.Int:
		if nullableInt, ok := field.Interface().(null.Int); ok {
			if nullableInt.Valid {
				return nullableInt.Int64
			}
			return 0
		}
	case null.Float:
		if nullableFloat, ok := field.Interface().(null.Float); ok {
			if nullableFloat.Valid {
				return nullableFloat.Float64
			}
			return 0
		}
	case null.Bool:
		if nullableBool, ok := field.Interface().(null.Bool); ok {
			if nullableBool.Valid {
				return nullableBool.Bool
			}
			return false
		}
	case null.Time:
		if nullableTime, ok := field.Interface().(null.Time); ok {
			if nullableTime.Valid {
				return nullableTime.Time
			}
			return time.Time{}
		}
	case nullable.Int:
		if nullableInt, ok := field.Interface().(nullable.Int); ok {
			if nullableInt.Valid {
				return nullableInt.Int64
			}
			return 0
		}
	case nullable.Float:
		if nullableFloat, ok := field.Interface().(nullable.Float); ok {
			if nullableFloat.Valid {
				return nullableFloat.Float64
			}
			return 0
		}
	case nullable.Bool:
		if nullableBool, ok := field.Interface().(nullable.Bool); ok {
			if nullableBool.Valid {
				return nullableBool.Bool
			}
			return false
		}
	case nullable.Time:
		if nullableTime, ok := field.Interface().(nullable.Time); ok {
			if nullableTime.Valid {
				return nullableTime.Time
			}
			return time.Time{}
		}
	case nullable.String:
		if nullableString, ok := field.Interface().(nullable.String); ok {
			if nullableString.Valid {
				return nullableString.String
			}
			return ""
		}
	}
	return nil
}

func ValidateAlphaNumericSpaceOnly(fl validator.FieldLevel) bool {
	// Regular expression for alphanumeric characters and spaces
	regexPattern := "^[a-zA-Z0-9 ]+$"

	// Compile the regular expression
	regex := regexp.MustCompile(regexPattern)

	// Use MatchString to check if the input matches the pattern
	return regex.MatchString(fl.Field().String())
}

func ValidateAlphaNumericSpaceOrEmpty(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	// Regular expression for alphanumeric characters and spaces
	regexPattern := "^[a-zA-Z0-9 ]+$"

	// Compile the regular expression
	regex := regexp.MustCompile(regexPattern)

	// Use MatchString to check if the input matches the pattern
	return regex.MatchString(fl.Field().String())
}

func ValidateNPWP(fl validator.FieldLevel) bool {
	// Format NPWP: 15 or 16 numeric digits only
	rp := `^\d{15,16}$`
	regex := regexp.MustCompile(rp)
	return regex.MatchString(fl.Field().String())
}

func ValidateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	format := fl.Param()
	if date == "" {
		return true
	}
	switch format {
	case "d-m-Y", "DD-MM-YYYY":
		format = "02-01-2006"
	case "Y-m-d", "YYYY-MM-DD":
		format = time.DateOnly
	default:
		format = "02-01-2006"
	}

	_, err := time.Parse(format, date)

	return err == nil
}

func ValidatePhoneNumber(fl validator.FieldLevel) bool {
	params := fl.Field().String()
	if params == "" || params == "-" {
		return true
	}
	rp := "^\\+?\\d+-?\\d+$"
	regex := regexp.MustCompile(rp)
	return regex.MatchString(fl.Field().String())
}

func ValidateInArray(fl validator.FieldLevel) bool {
	params := fl.Field().String()
	acceptedList := fl.Param()
	acceptedListArr := strings.Split(acceptedList, ";")

	return array_helper.InArray(acceptedListArr, params)
}

/*
	ValidateDecimal used to validate decimal precision and scale.

Format: decimal=precision;scale
Example: validate:decimal=20;4
*/
func ValidateDecimal(fl validator.FieldLevel) bool {
	switch fl.Field().Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return DecimalValidator(float64(fl.Field().Int()), fl.Param())
	case reflect.Float32, reflect.Float64:
		return DecimalValidator(fl.Field().Float(), fl.Param())
	case reflect.Struct:
		switch fl.Field().Type() {
		case reflect.TypeOf(null.Int{}), reflect.TypeOf(nullable.Int{}):
			return DecimalValidator(float64(fl.Field().Interface().(nullable.Int).Int64), fl.Param())
		case reflect.TypeOf(null.Float{}), reflect.TypeOf(nullable.Float{}):
			return DecimalValidator(fl.Field().Interface().(nullable.Float).Float64, fl.Param())
		}
	}

	return DecimalValidator(fl.Field().Float(), fl.Param())
}

func DecimalValidator(params float64, decimal string) bool {
	// Get precicion and scale for decimal
	decimalSplit := strings.Split(decimal, ";")
	if len(decimalSplit) != 2 {
		return false
	}
	precision, _ := strconv.Atoi(decimalSplit[0])
	scale, _ := strconv.Atoi(decimalSplit[1])

	parts := strings.Split(fmt.Sprintf("%."+strconv.Itoa(scale)+"f", params), ".")

	intPart := parts[0]
	fracPart := parts[1]

	// Check if the integer part length exceeds the allowed precision minus scale
	if len(intPart) > precision-scale {
		return false
	}

	// Check if the fractional part length exceeds the allowed scale
	if len(fracPart) > scale {
		return false
	}
	return true
}
