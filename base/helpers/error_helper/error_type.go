package error_helper

import (
	"core-ticket/constants/error_code"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"gopkg.in/guregu/null.v4"
)

/*
Error represents an error object that can contain multiple validation errors along with an error code, raw error message, and additional data.

Fields:
- errorCode: The main error code associated with the error.
- rawErrorMessage: The raw error message as a nullable string.
- validationErrors: A slice of ValidationError objects representing detailed validation errors.
- data: Additional data that can be attached to the error object.

Note: The validationErrors field can be used to store multiple ValidationError objects for detailed error reporting.
*/
type Error struct {
	errorCode        error_code.ErrorCode
	rawErrorMessage  null.String
	validationErrors []*ValidationError
	data             interface{}
}

/*
Error returns the error message based on the error code or the raw error message.

If the raw error message is valid, it returns the raw error message as a string.
If the raw error message is not valid, it returns the corresponding message for the error code.

Returns:
  - string: The error message.
*/
func (errs *Error) Error() string {
	if errs.rawErrorMessage.Valid {
		return errs.rawErrorMessage.String
	}
	return errs.errorCode.Message()
}

/*
Data returns the data associated with the error.
It can be used to retrieve additional information related to the error.
The data type is interface{}, allowing flexibility in the type of data stored.
*/
func (errs *Error) Data() interface{} {
	return errs.data
}

/*
Set updates the Error instance with the provided error information by calling the config method.
It sets the errorCode based on the type of error received and populates the validationErrors or rawErrorMessage accordingly.
If the error is of type sql.ErrNoRows, it sets the errorCode to error_code.NotFound.
If the error is of type validator.ValidationErrors, it sets the errorCode to error_code.ValidationError and populates the validationErrors slice with the individual validation errors.
If the error is of type *ValidationError, it sets the errorCode to error_code.ValidationError and appends the error to the validationErrors slice.
For any other type of error, it sets the errorCode to error_code.UnknownError and stores the error message in rawErrorMessage.
The updated Error instance is then returned.

Parameters:
- err: The error to be processed and updated in the Error instance.

Returns:
- *Error: The updated Error instance.
*/
func (errs *Error) Set(err error) *Error {
	return errs.config(err)
}

/*
config updates the Error instance based on the provided error.
If the error is sql.ErrNoRows, sets the errorCode to NotFound.
If the error is of type validator.ValidationErrors, sets the errorCode to ValidationError and populates the validationErrors slice with individual ValidationError instances.
If the error is a single ValidationError, adds it to the validationErrors slice.
If the error is of any other type, sets the errorCode to UnknownError and stores the error message in rawErrorMessage.

Parameters:
- err: The error to be processed.

Returns:
- The updated Error instance.
*/
func (errs *Error) config(err error) *Error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		errs.errorCode = error_code.NotFound
	case errors.As(err, &validator.ValidationErrors{}):
		errs.errorCode = error_code.ValidationError

		for _, fe := range err.(validator.ValidationErrors) {
			ve := &ValidationError{}
			errs.validationErrors = append(errs.validationErrors, ve.Set(fe))
		}
	default:
		var v *ValidationError
		if errors.As(err, &v) {
			errs.errorCode = error_code.ValidationError
			errs.validationErrors = append(errs.validationErrors, err.(*ValidationError))
		} else {
			errs.errorCode = error_code.UnknownError
			if err != nil {
				errs.rawErrorMessage = null.StringFrom(err.Error())
			} else {
				errs.rawErrorMessage.Valid = false
			}
		}
	}

	return errs
}

/*
SetValidationErrors sets the validation errors for the Error object.

Parameters:
- e: a slice of ValidationError objects representing the validation errors to be set.
- data: optional parameter to set additional data for the Error object.

Returns:
- *Error: the updated Error object with the validation errors set.

Example:
err := &Error{}
validationErr := &ValidationError{ErrorCode: error_code.ValidationError, Tag: "required", Field: "name", Param: "", Row: 1, message: "Name is required"}
err.SetValidationErrors([]*ValidationError{validationErr}, "additional data")
*/
func (errs *Error) SetValidationErrors(e []*ValidationError, data ...interface{}) *Error {
	errs.errorCode = error_code.ValidationError
	for _, err := range e {
		errs.validationErrors = append(errs.validationErrors, err)
	}

	if len(data) > 0 {
		errs.data = data[0]
	}

	return errs
}

/*
ErrorCode returns the error code associated with the Error instance.

Returns:

	error_code.ErrorCode: The error code of the Error instance.
*/
func (errs *Error) ErrorCode() error_code.ErrorCode {
	return errs.errorCode
}

/*
setErrorCode sets the error code for the Error instance.

Parameters:
- errorCode: a variadic parameter of type error_code.ErrorCode, representing the error code to be set.

Returns:
- *Error: the updated Error instance with the specified error code.

Example:

	err := &Error{}
	err = err.setErrorCode(error_code.ErrorCode("12345"))

Note:
If multiple error codes are provided, only the first one will be set.
*/
func (errs *Error) setErrorCode(errorCode ...error_code.ErrorCode) *Error {
	if errorCode != nil {
		errs.errorCode = errorCode[0]
	}
	return errs
}

/*
ValidationErrors returns a slice of validation errors associated with the error.
*/
func (errs *Error) ValidationErrors() []*ValidationError {
	return errs.validationErrors
}

/*
New creates a new Error instance with the provided error and error code.

Parameters:
- err: The error to be set in the Error instance.
- errorCode: Optional error code to be set in the Error instance.

Returns:
- *Error: The newly created Error instance.

Example:
err := errors.New(errors.New("Something went wrong"), error_code.InternalServerError)
*/
func New(err error, errorCode ...error_code.ErrorCode) *Error {
	return (&Error{}).Set(err).setErrorCode(errorCode...)
}

/*
IsNotFound checks if the given error is of type *Error and has an error code of NotFound (00004).
It returns true if the error is a NotFound error, otherwise it returns false.

Parameters:
- e: The error to be checked.

Returns:
- bool: True if the error is a NotFound error, false otherwise.
*/
func IsNotFound(e error) bool {
	var er *Error
	if errors.As(e, &er) && er.errorCode == error_code.NotFound {
		return true
	}
	return false
}

/*
Is checks if the given error is of type *Error and returns true if it is, otherwise returns false.

Parameters:
- e: The error to check.

Returns:
- bool: True if the error is of type *Error, false otherwise.
*/
func Is(e error) bool {
	var er *Error
	if errors.As(e, &er) {
		return true
	}

	return false
}

/*
As takes an error as input and attempts to convert it to an *Error type.
If successful, it returns the converted *Error, otherwise it returns nil.

Parameters:
- e: The error to be converted.

Returns:
- *Error: The converted *Error if successful, otherwise nil.
*/
func As(e error) *Error {
	var er *Error
	if errors.As(e, &er) {
		return er
	}

	return nil
}
