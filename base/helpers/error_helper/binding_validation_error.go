package error_helper

import (
	"core-ticket/constants/error_code"
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

const (
	LOOKUP_JSON = "json"
	LOOKUP_FORM = "form"
)

/*
ParseValidatorError parses the given error and generates a list of validation errors based on the lookup type and model provided.

Parameters:
- err: The error to be parsed.
- lookupType: The type of lookup to be used for validation.
- model: The model to be validated.

Returns:
- error: An error indicating the validation errors encountered during parsing.

Deprecated: This function is deprecated and should be replaced with a more efficient validation error handling mechanism.
*/
func ParseValidatorError(err error, lookupType string, model interface{}) error {
	var validationErrors []*ValidationError

	structType := reflect.TypeOf(model)
	if structType.Kind() != reflect.Ptr || structType.Elem().Kind() != reflect.Struct {
		return New(errors.New("invalid struct type"), error_code.UnknownError)
	}
	structType = structType.Elem()

	switch v := err.(type) {
	case validator.ValidationErrors:
		for _, validationErr := range v {
			fieldName := validationErr.Field()
			field, _ := structType.FieldByName(fieldName)
			fieldLookupName, _ := field.Tag.Lookup(lookupType)
			if fieldLookupName != "" {
				validationErrors = append(validationErrors, NewValidationError(
					error_code.ValidationErrorCode(error_code.ValidationError),
					validationErr.ActualTag(),
					fieldLookupName,
					"",
					0,
					validationErr.Error(),
				))
			} else {
				validationErrors = append(validationErrors, NewValidationError(
					error_code.ValidationErrorCode(error_code.ValidationError),
					validationErr.ActualTag(),
					fieldName,
					"",
					0,
					validationErr.Error(),
				))
			}
		}

	}
	return New(nil, error_code.ValidationError).SetValidationErrors(validationErrors)
}

/*
ParseValidatorErrorV2 parses the given error and generates a list of validation errors based on the lookup type and model provided.

Parameters:
- err: The error to be parsed.
- lookupType: The type of lookup to be used for validation.
- model: The model to be validated.

Returns:
- error: An error indicating the validation errors encountered during parsing.
*/
func ParseValidatorErrorV2(err error, lookupType string, model interface{}) error {
	var validationErrors []*ValidationError

	// Get model type
	modelValue := reflect.ValueOf(model)
	for modelValue.Kind() == reflect.Ptr || modelValue.Kind() == reflect.Interface {
		modelValue = modelValue.Elem()
	}
	// Check if the modelValue is a struct
	if modelValue.Kind() != reflect.Struct {
		return New(errors.New("invalid struct type"), error_code.UnknownError)
	}

	// Get model type
	modelType := modelValue.Type()

	// Set validation errors
	validationErrors = getFieldNested(err, lookupType, modelType, validationErrors)

	return New(nil, error_code.ValidationError).SetValidationErrors(validationErrors)
}

func getFieldNested(err error, lookupType string, model reflect.Type, validationErrors []*ValidationError) []*ValidationError {
	// Loop through each field in the struct
	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)

		// If field is a slice or array, get the element field
		if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			elemType := field.Type.Elem()
			if !isPrimitiveType(elemType) {
				validationErrors = getFieldNested(err, lookupType, elemType, validationErrors)
			}
		}
		// Check if the field is a struct
		if field.Type.Kind() == reflect.Struct {
			// If field is a struct,
			validationErrors = getFieldNested(err, lookupType, field.Type, validationErrors)
		} else {
			// Get tag name based on lookup type
			lookupTagName, _ := field.Tag.Lookup(lookupType)

			// Loop through errors
			switch v := err.(type) {
			case validator.ValidationErrors:
				for _, validationErr := range v {
					fieldName := validationErr.Field()
					// If field is same with the validator, set error
					if fieldName == field.Name {
						validationErrors = append(validationErrors, NewValidationError(
							error_code.ValidationErrorCode(error_code.ValidationError),
							validationErr.ActualTag(),
							lookupTagName,
							"",
							0,
							validationErr.Error(),
						))
						break
					}
				}
			}
		}
	}

	return validationErrors
}

func isPrimitiveType(t reflect.Type) bool {
	// Check if the type is one of the primitive data types
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Bool:
		return true
	case reflect.String:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Interface:
		return false
	default:
		return false
	}
}
