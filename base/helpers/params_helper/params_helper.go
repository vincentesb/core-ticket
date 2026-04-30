package params_helper

import (
	"reflect"
	"time"

	"github.com/oleiade/reflections"
	"gopkg.in/guregu/null.v4"
)

// SetDefaultParams used to set default parameters, including CreatedDate, CreatedBy, EditedBy, EditedDate
// It returns param that have been assigned for its default value
func SetDefaultParams[T interface{}](username string, param T, isCreate bool) T {
	now := time.Now()

	if isCreate {
		_ = reflections.SetField(&param, "CreatedBy", username)
		_ = reflections.SetField(&param, "CreatedDate", now)
	}

	_ = reflections.SetField(&param, "EditedBy", null.StringFrom(username))
	_ = reflections.SetField(&param, "EditedDate", null.TimeFrom(now))

	return param
}

// EnforceNilSlices iterates through the fields of any struct passed as input,
// and if it finds any slice fields that are empty, it sets them to nil.
func EnforceNilSlices(input interface{}) {
	v := reflect.ValueOf(input).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Slice && field.Len() == 0 {
			field.Set(reflect.Zero(field.Type()))
		}
	}
}
