package struct_helper

import (
	"fmt"
	"reflect"
	"strings"
)

/*
CloneStruct creates a deep copy of the provided struct 'source' and returns it.
It uses reflection to copy the values of 'source' to the new struct 'output'.
The function expects 'source' to be a struct and returns a new struct of the same type.

Parameters:
- source: The struct to be cloned.

Returns:
- T: A new struct that is a deep copy of 'source'.
*/
func CloneStruct[T interface{}](source T) T {
	var output T
	vSource := reflect.ValueOf(source)
	reflect.ValueOf(&output).Elem().Set(vSource)
	return output
}

/*
CreateStructWithCommonField takes an input source interface{} and creates a new struct of type T with common fields populated from the source.
It uses reflection to map fields from the source to the output struct based on field names or tags.
The function returns the newly created struct of type T with common fields populated from the source.

Parameters:
- source: an interface{} representing the input source from which common fields will be populated.

Returns:
- T: a new struct of type T with common fields populated from the source.
*/
func CreateStructWithCommonField[T interface{}](source interface{}) T {
	var output T
	rvSource := reflect.ValueOf(source)
	rvOutput := reflect.ValueOf(&output).Elem()

	internalConverter(rvSource, rvOutput)

	return output
}

/*
ConvertStructWithCommonField takes in a source interface{} and a destination interface{}.
It converts the fields from the source struct to the destination struct based on common field names.
The function uses reflection to achieve this conversion.
*/
func ConvertStructWithCommonField(source interface{}, destination interface{}) {
	rvSource := reflect.ValueOf(source)
	rvOutput := reflect.ValueOf(destination).Elem()

	internalConverter(rvSource, rvOutput)
}

/*
internalConverter takes in two reflect.Value arguments, rvSource and rvOutput, and iterates over the fields of rvOutput to map corresponding fields from rvSource.
It handles cases where fields are pointers, null types, or need type conversion.
The function uses reflection to achieve this mapping between the source and output structs.
*/
func internalConverter(rvSource reflect.Value, rvOutput reflect.Value) {
	for i := 0; i < rvOutput.NumField(); i++ {
		ftOutput := rvOutput.Type().Field(i)

		fOutputName := rvOutput.Type().Field(i).Name
		if value, ok := ftOutput.Tag.Lookup("map"); ok {
			fOutputName = value
		}
		fSource := rvSource.FieldByName(fOutputName)
		if !fSource.IsValid() {
			continue
		}

		switch fSource.Kind() {
		case reflect.Ptr:
			if fSource.IsNil() {
				continue
			}
			fSource = fSource.Elem()
		}

		fOutput := rvOutput.Field(i)
		if fOutput.Kind() == reflect.Ptr {
			if fOutput.IsNil() {
				fOutput.Set(reflect.New(fOutput.Type().Elem()))
			}
			fOutput = fOutput.Elem()
		}

		if strings.Contains(fSource.Type().PkgPath(), "null") {
			nullStruct := fSource.FieldByNameFunc(func(s string) bool {
				if strings.Contains(s, "Null") {
					return true
				}
				return false
			})
			// Check Valid flag (Field(1))
			if !nullStruct.Field(1).Bool() {
				// Source is NULL
				if strings.Contains(fOutput.Type().PkgPath(), "null") {
					// Destination is also null type, set it to NULL
					fOutputNull := fOutput.FieldByNameFunc(func(s string) bool {
						if strings.Contains(s, "Null") {
							return true
						}
						return false
					})
					fOutputNull.Field(1).SetBool(false) // Set Valid = false
					fOutputNull.Field(0).Set(reflect.Zero(fOutputNull.Field(0).Type()))
				}
				// Skip this field (can't set NULL to non-null type)
				continue
			}
			fSource = nullStruct.Field(0)
		}

		if strings.Contains(fOutput.Type().PkgPath(), "null") {
			fOutput = fOutput.FieldByNameFunc(func(s string) bool {
				if strings.Contains(s, "Null") {
					return true
				}
				return false
			})
			fOutput.Field(1).SetBool(true)
			fOutput = fOutput.Field(0)
		}

		if fOutput.Kind() == reflect.Int64 && fSource.Kind() == reflect.Int {
			fOutput.SetInt(fSource.Int())
		}

		if fOutput.Kind() == fSource.Kind() &&
			fOutput.Type() == fSource.Type() {
			fOutput.Set(fSource)
		}
	}
}

/*
GetDBTagFromStruct retrieves the 'db' tags from the fields of the provided struct or pointer to a struct.
It iterates over the fields of the struct, checks for the presence of 'db' tag, and appends the tag value to the output slice.
If the provided interface is a pointer to a struct or a slice of structs, it dereferences the pointer or extracts the element type of the slice respectively.

Parameters:
- dest: The struct or pointer to a struct from which 'db' tags will be extracted.

Returns:
- []string: A slice containing the 'db' tags found in the fields of the struct.

Example:

	dbTags := GetDBTagFromStruct(&MyStruct{})
	for _, tag := range dbTags {
		fmt.Println(tag)
	}
*/
func GetDBTagFromStruct(dest interface{}, alias ...string) []string {
	rt := reflect.TypeOf(dest)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	var dbTags []string
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		dbTag := f.Tag.Get("db")
		if dbTag == "" {
			continue
		}
		if alias != nil {
			dbTags = append(dbTags, fmt.Sprintf("`%s`.`%s`", alias[0], dbTag))
		} else {
			dbTags = append(dbTags, fmt.Sprintf("`%s`", dbTag))
		}
	}

	return dbTags
}

// IsEmpty used to validate neither struct is empty
func IsEmpty(data interface{}) bool {
	val := reflect.ValueOf(data)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		zero := reflect.Zero(field.Type())

		if !reflect.DeepEqual(field.Interface(), zero.Interface()) {
			return false
		}
	}

	return true
}
