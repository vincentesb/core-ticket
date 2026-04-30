package base_helper

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

const (
	SCENARIO_CREATE    = "create"
	SCENARIO_UPDATE    = "update"
	SCENARIO_DELETE    = "delete"
	COMPARE_GET_BEFORE = 1
	COMPARE_GET_AFTER  = 2
)

func ConvertToInteger(variable interface{}) int {
	var number int

	switch val := variable.(type) {
	case int:
		number = val
	case float64:
		number = int(val)
	case string:
		number, _ = strconv.Atoi(val)
	}

	return number
}

func ConvertToBoolean(variable interface{}) (bool, error) {
	var boolean bool
	switch val := variable.(type) {
	case int:
		if val > 0 {
			boolean = true
		} else {
			boolean = false
		}
	case bool:
		boolean = val
	case string:
		boolean, _ = strconv.ParseBool(val)
	}
	return boolean, nil
}

func ConvertToFloat(variable interface{}) float64 {
	var float float64
	switch val := variable.(type) {
	case int:
		float = float64(val)
	case float64:
		float = val
	case string:
		float, _ = strconv.ParseFloat(val, 64)
	}
	return float
}

func ConvertStructWithCommonField(source interface{}, destination interface{}) {
	dst := reflect.ValueOf(destination).Elem()
	src := reflect.ValueOf(source).Elem()

	for i := 0; i < dst.Type().NumField(); i++ {
		dstField := dst.Type().Field(i)

		var srcFieldValue reflect.Value
		if mapTag, isValid := dstField.Tag.Lookup("map"); isValid {
			srcFieldValue = src.FieldByName(mapTag)
		} else {
			srcFieldValue = src.FieldByName(dstField.Name)
		}

		if !srcFieldValue.IsValid() {
			continue
		}

		dstFieldValue := dst.Field(i)
		if dstFieldValue.Kind() == reflect.Ptr {
			if dstFieldValue.IsNil() {
				dstFieldValue.Set(reflect.New(dstFieldValue.Type().Elem()))
			}
			dstFieldValue = dstFieldValue.Elem()
		}
		if srcFieldValue.Kind() == reflect.Ptr {
			if srcFieldValue.IsNil() {
				srcFieldValue.Set(reflect.New(srcFieldValue.Type().Elem()))
			}
			srcFieldValue = srcFieldValue.Elem()
		}

		dstFieldType := dstFieldValue.Type().String()

		if strings.Contains(dstFieldType, TYPE_INTEGER) {
			dstFieldValue.Set(reflect.ValueOf(Integer(srcFieldValue.Int())))
		} else if strings.Contains(dstFieldType, TYPE_BOOLEAN) {
			dstFieldValue.Set(reflect.ValueOf(Boolean(srcFieldValue.Bool())))
		} else if strings.Contains(dstFieldType, TYPE_FLOAT) {
			dstFieldValue.Set(reflect.ValueOf(Float(srcFieldValue.Float())))
		} else {
			switch dstFieldValue.Kind() {
			case reflect.Int:
				dstFieldValue.SetInt(srcFieldValue.Int())
			case reflect.Bool:
				dstFieldValue.SetBool(srcFieldValue.Bool())
			case reflect.String:
				dstFieldValue.SetString(srcFieldValue.String())
			case reflect.Float64:
				dstFieldValue.SetFloat(srcFieldValue.Float())
			default:
				dstFieldValue.Set(srcFieldValue)
			}
		}
	}
}

func CompareStruct(a interface{}, b interface{}, option int) map[string]interface{} {
	differingKeysValues := make(map[string]interface{})

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() == reflect.Slice {
		if va.Kind() == reflect.Slice && vb.Kind() == reflect.Slice {
			differences := compareSlices(va, vb, option)
			if len(differences) > 0 {
				differingKeysValues["changes"] = differences
			}
		}
	} else {
		for i := 0; i < va.NumField(); i++ {
			field := va.Type().Field(i)
			valA := va.Field(i)
			valB := vb.Field(i)

			if !valA.CanInterface() || !valB.CanInterface() {
				continue
			}

			if !reflect.DeepEqual(valA.Interface(), valB.Interface()) {
				if valA.Kind() == reflect.Slice && valB.Kind() == reflect.Slice {
					differences := compareSlices(valA, valB, option)
					if len(differences) > 0 {
						differingKeysValues[field.Name] = differences
					}
				} else {
					if option == COMPARE_GET_BEFORE {
						differingKeysValues[field.Name] = valA.Interface()
					} else {
						differingKeysValues[field.Name] = valB.Interface()
					}
				}
			}
		}
	}

	return differingKeysValues
}

func compareSlices(sliceA, sliceB reflect.Value, option int) []interface{} {
	var differences []interface{}

	lenA := sliceA.Len()
	lenB := sliceB.Len()
	maxLen := lenA
	if lenB > lenA {
		maxLen = lenB
	}

	for i := 0; i < maxLen; i++ {
		var elementA, elementB interface{}
		if i < lenA {
			elementA = sliceA.Index(i).Interface()
		}
		if i < lenB {
			elementB = sliceB.Index(i).Interface()
		}

		if reflect.DeepEqual(elementA, elementB) {
			// Elements are equal, no need to append
			continue
		}

		if option == COMPARE_GET_BEFORE {
			differences = append(differences, elementA)
		} else {
			differences = append(differences, elementB)
		}
	}

	return differences
}

func CloneStruct[T interface{}](source T) T {
	vSource := reflect.ValueOf(source)
	var res T
	vRes := reflect.ValueOf(&res).Elem()

	for i := 0; i < vSource.NumField(); i++ {
		switch vSource.Field(i).Kind() {
		case reflect.Ptr:
			if vRes.Field(i).IsNil() {
				//vRes.Field(i).Set(reflect.New(vRes.Field(i).Type().Elem()))
				continue
			}

			switch vSource.Field(i).Elem().Kind() {
			case reflect.Int:
				vRes.Field(i).Elem().SetInt(vSource.Field(i).Elem().Int())
			case reflect.String:
				vRes.Field(i).Elem().SetString(vSource.Field(i).Elem().String())
			case reflect.Bool:
				vRes.Field(i).Elem().SetBool(vSource.Field(i).Elem().Bool())
			default:
				vRes.Field(i).Elem().Set(vSource.Field(i).Elem())
			}
		case reflect.Int:
			vRes.Field(i).SetInt(vSource.Field(i).Int())
		case reflect.Bool:
			vRes.Field(i).SetBool(vSource.Field(i).Bool())
		case reflect.String:
			vRes.Field(i).SetString(vSource.Field(i).String())
		default:
			vRes.Field(i).Set(vSource.Field(i))
		}
	}

	return res
}

func GetMapKeys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetMapValues[T comparable, U any](m map[T]U) []U {
	values := make([]U, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Empty(variable interface{}) bool {
	rV := reflect.ValueOf(variable)
	if rV.Kind() == reflect.Ptr {
		rV = rV.Elem()
	}
	if !rV.IsValid() || rV.IsZero() {
		return true
	}
	return false
}

func GetAddress[T interface{}](variable T) *T {
	if Empty(variable) {
		return nil
	}
	return &variable
}

func CreateNullString(variable string) null.String {
	return null.StringFromPtr(GetAddress(variable))
}

func GetTotal(qty float64, price float64, discount float64, vatValue float64, taxRate float64) float64 {
	beforeVatSubtotal := (price * qty) - discount
	vatTotal := beforeVatSubtotal * vatValue / 100
	taxTotal := beforeVatSubtotal * taxRate / 100

	return beforeVatSubtotal + vatTotal - taxTotal
}

func ContainsComplexCharacters(input string) bool {
	complexCharRegex := regexp.MustCompile(`[一-龯㐀-䶵ぁ-ゟ゠-ヿㄱ-ㅣ가-힣⺀-⻳⼀-⿕㌀-㏿㐀-䶵阿-鿿𠀀-𫜴؀-ۿ܀-ݿހ-޿]`)
	return complexCharRegex.MatchString(input)
}
