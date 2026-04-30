package query_helper

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
Deprecated: This struct is deprecated and should not be used. Use query_builder.QueryBuilder instead
*/
type WhereCondition struct {
	Operator string
	Value    interface{}
}

/*
Deprecated: This struct is deprecated and should not be used. Use query_builder.QueryBuilder instead
*/
type SortCondition struct {
	Direction string
}

/*
Deprecated: This struct is deprecated and should not be used. Use query_builder.QueryBuilder instead
*/
type LimitCondition struct {
	Skip int
	Take int
}

/*
GenerateConditionQuery generates a query based on the provided conditions and returns the parameters to be used in the query.

Parameters:
- baseQuery: A pointer to the base query string to which conditions will be added.
- condition: A map containing the conditions to be applied to the query. The key represents the column name, and the value is a WhereCondition struct containing the operator and value.

Returns:
- []interface{}: A slice of interface{} containing the parameters to be used in the query.
- error: An error if any validation fails during the query generation process.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Note: The function dynamically constructs a query based on the conditions provided in the map. It handles operators like IN, LIKE, and BETWEEN to build the query string and parameters accordingly.
*/
func GenerateConditionQuery(
	baseQuery *string,
	condition map[string]WhereCondition,
) ([]interface{}, error) {
	var params []interface{}
	for key, val := range condition {
		rfVal := reflect.ValueOf(val.Value)
		if rfVal.Kind() == reflect.Ptr {
			rfVal = rfVal.Elem()
		}
		if !rfVal.IsValid() {
			continue
		}

		if key != "" && val.Operator != "" {
			switch {
			case strings.Contains(val.Operator, "IN"):
				if rfVal.Kind() != reflect.Slice {
					return nil, errors.New("value type of IN operator must be Slice")
				}
				if rfVal.IsNil() {
					continue
				}
				*baseQuery += " AND " + key
				*baseQuery += " " + val.Operator + " ("
				for i := 0; i < rfVal.Len(); i++ {
					*baseQuery += "?"
					if i != rfVal.Len()-1 {
						*baseQuery += ","
					}
					params = append(params, rfVal.Index(i).Interface())
				}
				*baseQuery += ")"
			case strings.Contains(val.Operator, "LIKE"):
				if rfVal.String() == "" {
					continue
				}
				if rfVal.Kind() != reflect.String {
					return nil, errors.New("value type of LIKE operator must be String")
				}
				*baseQuery += " AND " + key
				*baseQuery += " " + val.Operator + " ?"
				params = append(params, "%"+rfVal.String()+"%")
			case strings.Contains(val.Operator, "BETWEEN"):
				if rfVal.Kind() != reflect.Slice {
					return nil, errors.New("value type of BETWEEN operator must be Slice")
				}
				if rfVal.Len() != 2 {
					return nil, errors.New("slice value must consist of two values")
				}
				value1 := rfVal.Index(0)
				value2 := rfVal.Index(1)
				if value1.Kind() == reflect.Interface && value1.Elem().Kind() == reflect.Ptr {
					value1 = value1.Elem()
				}
				if value2.Kind() == reflect.Interface && value2.Elem().Kind() == reflect.Ptr {
					value2 = value2.Elem()
				}
				if value1.IsNil() || value2.IsNil() {
					continue
				}
				*baseQuery += " AND " + key
				*baseQuery += " " + val.Operator + " ? AND ?"
				params = append(params, value1.Interface(), value2.Interface())
			default:
				if rfVal.String() == "" || rfVal.IsZero() {
					continue
				}
				*baseQuery += " AND " + key
				*baseQuery += " " + val.Operator + " ?"
				params = append(params, rfVal.Interface())
			}
		}
	}
	return params, nil
}

/*
GenerateSortConditionQuery generates a SQL query string based on the provided base query and sort conditions.

Parameters:
- baseQuery: a pointer to a string representing the base SQL query.
- sort: a map containing field names as keys and SortCondition structs as values.

SortCondition struct:
- Direction: a string representing the sorting direction (e.g., "ASC" or "DESC").

The function iterates over the sort map, checks if the field exists using the IsFieldExist function, and appends the sorting condition to the base query string.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Example:

	baseQuery := "SELECT * FROM table"
	sort := map[string]SortCondition{
		"field1": {Direction: "ASC"},
		"field2": {Direction: "DESC"},
	}
	GenerateSortConditionQuery(&baseQuery, sort)

This will modify baseQuery to "SELECT * FROM table ORDER BY field1 ASC, field2 DESC".
*/
func GenerateSortConditionQuery[T interface{}](
	baseQuery *string,
	sort map[string]SortCondition,
) {
	idx := 0
	for key, val := range sort {
		if key != "" {
			if idx == 0 {
				*baseQuery += " ORDER BY "
			} else {
				*baseQuery += ", "
			}
			if IsFieldExist[T](key) {
				*baseQuery += key + " " + val.Direction
			}
			idx++
		}
	}
}

/*
GenerateLimitConditionQuery generates a SQL query string with a LIMIT clause based on the provided LimitCondition struct.

Parameters:
- baseQuery: A pointer to a string representing the base SQL query.
- limit: A pointer to a LimitCondition struct containing the Skip and Take values for the LIMIT clause.

Deprecated: This function is deprecated and should not be used. Use query_builder.QueryBuilder instead

Example:

	GenerateLimitConditionQuery(&baseQuery, &LimitCondition{Skip: 10, Take: 5})

This function appends a LIMIT clause to the baseQuery string using the Skip and Take values from the limit struct. If the limit parameter is nil, no LIMIT clause is added to the query.
*/
func GenerateLimitConditionQuery(baseQuery *string, limit *LimitCondition) {
	if limit != nil {
		*baseQuery += " LIMIT " + strconv.Itoa(limit.Skip) + "," + strconv.Itoa(limit.Take)
	}
}

/*
GenerateRawInConditionValue generates a comma-separated list of values for an IN condition based on the provided slice.

Parameters:
- slice: An interface representing the slice of values to be included in the IN condition.

Deprecated: This function is deprecated and should not be used. Because it is not safe to use (e.g. sql injection),
Use GenerateInConditionPlaceholdersArgs instead

Returns:
- A string containing the comma-separated list of values for the IN condition.
*/
func GenerateRawInConditionValue(slice interface{}) string {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return ""
	}

	conditions := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		switch val := item.(type) {
		case string:
			conditions[i] = fmt.Sprintf("'%s'", val)
		case int, int8, int16, int32, int64:
			conditions[i] = fmt.Sprintf("%d", val)
		case bool:
			conditions[i] = fmt.Sprintf("%t", val)
		default:
			conditions[i] = fmt.Sprintf("'%v'", val)
		}
	}
	return strings.Join(conditions, ", ")
}

type InConditionType interface {
	string | bool |
		int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		time.Time
}

// GenerateInConditionPlaceholdersArgs returns a comma-separated list of placeholders
// and a slice of interface{} arguments for safe SQL query usage (e.g., IN (?, ?, ?)).
//
// If the input slice is empty, it returns a single placeholder and the zero value of the type.
//
// Example:
//
//	ph, args := GenerateInConditionPlaceholdersArgs([]string{"A", "B"})
//	// ph: "?, ?"
//	// args: []interface{}{"A", "B"}
func GenerateInConditionPlaceholdersArgs[T InConditionType](slice []T) (string, []interface{}) {
	if len(slice) == 0 {
		var zero T
		return "?", []interface{}{zero}
	}

	placeholders := make([]string, len(slice))
	args := make([]interface{}, len(slice))

	for i, val := range slice {
		placeholders[i] = "?"
		args[i] = val
	}

	return strings.Join(placeholders, ", "), args
}
