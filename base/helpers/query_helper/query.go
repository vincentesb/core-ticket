package query_helper

import (
	"reflect"
	"regexp"
	"strings"
)

/*
IsFieldExist checks if a field with the specified tag exists in the struct type T.
It takes the fieldInput string as the tag to search for.
Returns true if the field with the specified tag exists, otherwise returns false.

Parameters:
- fieldInput: string - The tag to search for in the struct fields.

Returns:
- bool: true if the field with the specified tag exists, false otherwise.
*/
func IsFieldExist[T interface{}](fieldInput string) bool {
	ctx := reflect.TypeOf((*T)(nil)).Elem()
	for i := 0; i < ctx.NumField(); i++ {
		field := ctx.Field(i)
		dbTag := field.Tag.Get("db")

		if strings.Compare(dbTag, fieldInput) == 0 {
			return true
		}
	}
	return false
}

/*
GetSort takes a string parameter 'sort' and returns two strings: 'sortColumn' and 'sortDirection'.
If the 'sort' string starts with "-", the 'sortDirection' is set to "DESC", otherwise it is set to "ASC".
The 'sortColumn' is set to 'sort' without the leading "-" character if present, otherwise it is set to 'sort'.

Parameters:
- sort: string - The input string representing the sort order.

Returns:
- sortColumn: string - The column to sort on.
- sortDirection: string - The direction of sorting, either "ASC" or "DESC".
*/
func GetSort(sort string) (sortColumn string, sortDirection string) {
	sortDirection = "ASC"
	//if sort has "-" the direction is descending
	if strings.HasPrefix(sort, "-") {
		sortDirection = "DESC"
		//remove the "-" character
		sortColumn = sort[:0] + sort[1:]
	} else {
		sortColumn = sort
	}

	return sortColumn, sortDirection
}

/*
GetSkip calculates the number of items to skip based on the given page number and limit.
It multiplies the page number by the limit and then subtracts the limit from the result.

Parameters:
- page: int - The page number for which to calculate the skip.
- limit: int - The maximum number of items per page.

Returns:
- skip: int - The number of items to skip based on the page and limit.
*/
func GetSkip(page int, limit int) (skip int) {
	skip = (page * limit) - limit
	return
}

/*
GetSkipTake calculates the number of items to skip and the number of items to take based on the given page number and limit.
It calculates the skip value by subtracting 1 from the page number and then multiplying it by the limit.
The take value is simply set to the provided limit.

Parameters:
- page: int - The page number for which to calculate the skip and take values.
- limit: int - The maximum number of items to take per page.

Returns:
- skip: int - The number of items to skip based on the page and limit.
- take: int - The number of items to take, which is equal to the limit.

Deprecated: This function is deprecated and should not be used. Use GetSkip
*/
func GetSkipTake(page int, limit int) (skip int, take int) {
	skip = (page - 1) * limit
	take = limit
	return
}

/*
IsFieldExistByTag checks if a field with the specified tag exists in the struct type T.
It takes the fieldInput string as the tag to search for.
Returns true if the field with the specified tag exists, otherwise returns false.

Parameters:
- fieldInput: string - The tag to search for in the struct fields.

Returns:
- bool: true if the field with the specified tag exists, false otherwise.
*/
func IsFieldExistByTag[T interface{}](fieldInput string) bool {
	ctx := reflect.TypeOf((*T)(nil)).Elem()
	for i := 0; i < ctx.NumField(); i++ {
		field := ctx.Field(i)
		dbTag := field.Tag.Get("sort")

		if strings.Compare(dbTag, fieldInput) == 0 {
			return true
		}
	}
	return false
}

// SanitizeSQLValue removes potentially dangerous SQL content from the input string
// to help mitigate SQL injection risks.
//
// This function performs the following steps:
//  1. Removes SQL comment patterns such as "--" and "#" and any content following them.
//  2. Removes dangerous SQL keywords (case-insensitive), including:
//     DROP, DELETE, INSERT, UPDATE, SELECT, UNION, ALTER, TRUNCATE, CREATE, REPLACE, EXEC, EXECUTE.
//  3. Removes special characters commonly used in SQL injection:
//     semicolon (;), single quote ('), double quote ("), and backslash (\).
//  4. Trims leading and trailing whitespace.
//
// Note: This function is meant as an additional safety layer for filtering user input,
// but it is **not a substitute** for parameterized or prepared SQL statements.
// Always use parameterized queries for real SQL query execution.
//
// Example:
//
//	input := "DROP TABLE users; --"
//	safe := SanitizeSQLValue(input)
//	fmt.Println(safe) // Output: "TABLE users"
func SanitizeSQLValue(input string) string {
	// Remove SQL comment patterns like -- and #
	commentPattern := regexp.MustCompile(`(?i)(--|#).*`)
	input = commentPattern.ReplaceAllString(input, "")

	// Remove SQL command keywords
	keywordPattern := regexp.MustCompile(`(?i)\b(DROP|DELETE|INSERT|UPDATE|SELECT|UNION|ALTER|TRUNCATE|CREATE|REPLACE|EXEC|EXECUTE)\b`)
	input = keywordPattern.ReplaceAllString(input, "")

	// Remove special characters
	dangerChars := regexp.MustCompile(`[;'"\\]`)
	input = dangerChars.ReplaceAllString(input, "")

	// Trim whitespace
	return strings.TrimSpace(input)
}
