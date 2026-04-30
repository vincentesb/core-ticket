package utility

import (
	"core-ticket/base/helpers/error_helper"
	cBinding "core-ticket/base/helpers/gin_helper/binding"
	"core-ticket/constants/error_code"
	"net/url"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidateStruct(obj any) error {
	validate := binding.Validator.Engine().(*validator.Validate)
	if err := validate.Struct(obj); err != nil {
		return error_helper.New(err, error_code.ValidationError)
	}
	return nil
}

func ValidateMap(values map[string]interface{}, rules map[string]interface{}) error {
	validate := binding.Validator.Engine().(*validator.Validate)
	errs := validate.ValidateMap(values, rules)

	var vErrors []*error_helper.ValidationError
	for k, v := range errs {
		for _, fe := range v.(validator.ValidationErrors) {
			vErrors = append(
				vErrors,
				error_helper.NewValidationError(
					error_code.ValidationErrorCode(error_code.ValidationError),
					fe.Tag(),
					k,
					fe.Param(),
					0,
				),
			)
		}
	}
	return error_helper.New(nil, error_code.ValidationError).SetValidationErrors(vErrors)
}

func GetRulesFromStruct(obj any, typ cBinding.Type) map[string]interface{} {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	rules := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		s := f.Tag.Get(typ.RulesTag())
		if s == "" {
			continue
		}
		key := strings.Split(s, ",")[0]
		rules[key] = f.Tag.Get("binding")
	}

	return rules
}

func ConvertUrlValuesToMapStringInterface(values url.Values) map[string]interface{} {
	res := make(map[string]interface{}, len(values))
	for k, value := range values {
		if len(value) > 1 {
			var temp []interface{}
			for _, v := range value {
				temp = append(temp, v)
			}
			res[k] = temp
		} else {
			res[k] = value[0]
		}
	}
	return res
}
