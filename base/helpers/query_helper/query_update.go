package query_helper

import "reflect"

type UpdateConfig struct {
	Model          interface{}
	SavedAttribute []string
}

/*
GenerateUpdateQuery generates an update query based on the provided UpdateConfig and appends the set clause to the baseQuery string pointer. It iterates over the SavedAttribute list of the UpdateConfig, retrieves the corresponding field value from the Model, and constructs the set clause for the update query. The function returns a slice of interface{} containing the values of the fields to be updated.

Parameters:
- updateConfig: UpdateConfig struct containing the model and saved attributes information.
- baseQuery: A pointer to a string representing the base query to which the set clause will be appended.

Returns:
- []interface{}: A slice of interface{} containing the values of the fields to be updated in the query.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Example:

	updateConfig := UpdateConfig{
		Model:          &User{ID: 1, Name: "Alice", Age: 30},
		SavedAttribute: []string{"Name", "Age"},
	}
	baseQuery := "UPDATE users"
	args := GenerateUpdateQuery(updateConfig, &baseQuery)
*/
func GenerateUpdateQuery(
	updateConfig UpdateConfig,
	baseQuery *string,
) []interface{} {
	var args []interface{}
	*baseQuery += " SET "
	for j, attr := range updateConfig.SavedAttribute {
		if j != 0 {
			*baseQuery += ","
		}
		var field reflect.StructField
		var value interface{}
		if reflect.TypeOf(updateConfig.Model).Kind() == reflect.Ptr {
			field, _ = reflect.TypeOf(updateConfig.Model).Elem().FieldByName(attr)
			value = reflect.ValueOf(updateConfig.Model).Elem().FieldByName(attr).Interface()
		} else {
			field, _ = reflect.TypeOf(updateConfig.Model).FieldByName(attr)
			value = reflect.ValueOf(updateConfig.Model).FieldByName(attr).Interface()
		}
		args = append(args, value)
		*baseQuery += field.Tag.Get("db") + " = ?"
	}

	return args
}
