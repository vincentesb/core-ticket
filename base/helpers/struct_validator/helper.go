package struct_validator

import (
	"core-ticket/base/helpers/error_helper"
	"core-ticket/constants/error_code"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lovelydeng/gomoji"
)

var timeFormat = "2006-01-02"

func GteDateBetween() validator.Func {
	return func(fl validator.FieldLevel) bool {
		target := fl.Param()
		targetVal, _ := fl.Parent().FieldByName(target).Interface().(string)
		sourceVal, _ := fl.Field().Interface().(string)
		sourceDate, _ := time.Parse(timeFormat, sourceVal)
		targetDate, _ := time.Parse(timeFormat, targetVal)
		return sourceDate.After(targetDate) || sourceDate.Equal(targetDate)
	}
}

func LteDateBetween() validator.Func {
	return func(fl validator.FieldLevel) bool {
		target := fl.Param()
		targetVal, _ := fl.Parent().FieldByName(target).Interface().(string)
		sourceVal, _ := fl.Field().Interface().(string)
		sourceDate, _ := time.Parse(timeFormat, sourceVal)
		targetDate, _ := time.Parse(timeFormat, targetVal)
		return sourceDate.Before(targetDate) || sourceDate.Equal(targetDate)
	}
}

func ValidateSort(sort string) error {
	validate := validator.New()
	if err := validate.Var(strings.TrimPrefix(sort, "-"), "alpha"); err != nil {
		var e validator.ValidationErrors
		if errors.As(err, &e) {
			fe := e[0]
			ve := error_helper.NewValidationError(
				error_code.IsNumeric,
				fe.Tag(),
				"sort",
				"",
				0,
			)

			return error_helper.New(
				nil,
				error_code.ValidationError,
			).SetValidationErrors([]*error_helper.ValidationError{ve})
		}
	}
	return nil
}

func NotEmoji() validator.Func {
	return func(fl validator.FieldLevel) bool {
		return !gomoji.ContainsEmoji(fl.Field().String())
	}
}
