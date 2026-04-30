package query_helper

import "strings"

/*
Deprecated: This struct is deprecated and should not be used. Use query_builder.QueryBuilder instead
*/
type SelectConfig struct {
	SelectedAttribute []string
}

/*
GenerateSelectQuery generates a select query based on the provided SelectConfig and baseQuery.
If the SelectConfig has SelectedAttribute, it returns nil. Otherwise, it constructs a query by replacing '*' in the baseQuery with placeholders for each selected attribute.

Parameters:
- selectConfig: SelectConfig struct containing the list of selected attributes.
- baseQuery: Pointer to a string representing the base query.

Returns:
- An array of interfaces representing the selected attributes.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Example:

	selectConfig := SelectConfig{
		SelectedAttribute: []string{"attr1", "attr2"},
	}
	baseQuery := "SELECT * FROM table"
	attrs := GenerateSelectQuery(selectConfig, &baseQuery)
	// Output: []interface{}{"attr1", "attr2"}
*/
func GenerateSelectQuery(
	selectConfig SelectConfig,
	baseQuery *string) []interface{} {
	if len(selectConfig.SelectedAttribute) > 0 {
		return nil
	}
	var attrs []interface{}
	var tempArgs *string
	for _, attr := range selectConfig.SelectedAttribute {
		*tempArgs += `?, `
		attrs = append(attrs, attr)
	}
	strings.Replace(*baseQuery, "*", *tempArgs, 1)
	return attrs
}
