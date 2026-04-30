package query_helper

import (
	"reflect"
)

type InsertConfig struct {
	Model          interface{}
	SavedAttribute []string
}

/*
GenerateInsertQuery generates an SQL insert query based on the provided InsertConfig and appends the query to the baseQuery string pointer. It returns a slice of interface{} containing the values to be inserted.

Parameters:
- insertConfig: InsertConfig struct containing the model and saved attributes for insertion.
- baseQuery: A pointer to a string where the generated query will be appended.

Returns:
- []interface{}: A slice of interface{} containing the values to be inserted into the database.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Note:
- This function uses reflection to dynamically generate the insert query based on the model and saved attributes provided in the InsertConfig struct.
*/
func GenerateInsertQuery(
	insertConfig InsertConfig,
	baseQuery *string,
) []interface{} {
	v := reflect.ValueOf(insertConfig.Model)
	var args []interface{}
	if v.Kind() == reflect.Slice {
		var sqlVarPlaceholder string
		for i := 0; i < v.Len(); i++ {
			model := v.Index(i).Interface()
			if i != 0 {
				*baseQuery += ","
			}
			if i == 0 {
				*baseQuery += "("
				sqlVarPlaceholder += "("
			}
			for j, attr := range insertConfig.SavedAttribute {
				if i == 0 {
					if j != 0 {
						*baseQuery += ","
						sqlVarPlaceholder += ","
					}
					field, _ := reflect.TypeOf(model).FieldByName(attr)
					*baseQuery += field.Tag.Get("db")
					sqlVarPlaceholder += "?"
					if j == len(insertConfig.SavedAttribute)-1 {
						*baseQuery += ") VALUES "
						sqlVarPlaceholder += ")"
					}
				}
				args = append(args, reflect.ValueOf(model).FieldByName(attr).Interface())
			}
			*baseQuery += sqlVarPlaceholder
		}
	} else {
		var sqlVarPlaceholder string
		*baseQuery += "("
		sqlVarPlaceholder += "("
		for j, attr := range insertConfig.SavedAttribute {
			if j != 0 {
				*baseQuery += ","
				sqlVarPlaceholder += ","
			}
			field, _ := reflect.TypeOf(insertConfig.Model).FieldByName(attr)
			*baseQuery += field.Tag.Get("db")
			sqlVarPlaceholder += "?"
			if j == len(insertConfig.SavedAttribute)-1 {
				*baseQuery += ") VALUES "
				sqlVarPlaceholder += ")"
			}
			args = append(args, reflect.ValueOf(insertConfig.Model).FieldByName(attr).Interface())
		}
		*baseQuery += sqlVarPlaceholder
	}
	return args
}
