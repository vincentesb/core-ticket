package array_helper

import (
	"fmt"
	"reflect"
)

// InArray used to check if value exists in array
// It returns true if exists, and false if not exists
func InArray[T comparable](src []T, val T) bool {
	for _, source := range src {
		if source == val {
			return true
		}
	}
	return false
}

// ArrayDiff used to get different value between two arrays
// It returns different value in array
func ArrayDiff[T comparable](arr1 []T, arr2 []T) []T {
	diff := make([]T, 0)
	for _, v1 := range arr1 {
		flagExists := false
		for _, v2 := range arr2 {
			if v1 == v2 {
				flagExists = true
				break
			}
		}
		if !flagExists {
			diff = append(diff, v1)
		}
	}

	for _, v1 := range arr2 {
		flagExists := false
		for _, v2 := range arr1 {
			if v1 == v2 {
				flagExists = true
				break
			}
		}
		if !flagExists {
			diff = append(diff, v1)
		}
	}

	return diff
}

// ArrayIntersectKey used to get intersect value from both array
// It returns map value
func ArrayIntersectKey[T1 comparable, T2 comparable](map1, map2 map[T1]T2) map[T1]T2 {
	intersect := make(map[T1]T2)
	for key := range map1 {
		if _, exists := map2[key]; exists {
			intersect[key] = map1[key]
		}
	}
	return intersect
}

// ArrayIntersect finds the intersection of two slices of any comparable type.
// It returns slice
func ArrayIntersect[T comparable](slice1, slice2 []T) []T {
	elementMap := make(map[T]bool)
	for _, elem := range slice1 {
		elementMap[elem] = true
	}

	var intersection []T
	for _, elem := range slice2 {
		if elementMap[elem] {
			intersection = append(intersection, elem)
		}
	}

	return intersection
}

// ConvertArrayToString used to convert array to string separated by separator
// It returns string value
func ConvertArrayToString[T1 comparable](arr []T1, separator string) string {
	stringValue := ""
	for i, val := range arr {
		if i == len(arr)-1 {
			stringValue += fmt.Sprint(val)
		} else {
			stringValue += fmt.Sprint(val) + separator
		}
	}
	return stringValue
}

// SliceToMap converts a slice of structs into a map based on the specified field name.
// The field name should be of type int or string.
func SliceToMap(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice but got %s", v.Kind())
	}

	result := make(map[interface{}]interface{})
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		itemValue := reflect.Indirect(item)
		key := itemValue.FieldByName(keyField)
		if !key.IsValid() {
			return nil, fmt.Errorf("field %s not found", keyField)
		}
		if !key.CanInterface() {
			return nil, fmt.Errorf("field %s cannot be interfaced", keyField)
		}

		keyInterface := key.Interface()
		switch key.Kind() {
		case reflect.Int, reflect.String:
			result[keyInterface] = item.Interface()
		default:
			return nil, fmt.Errorf("field %s must be of type int or string", keyField)
		}
	}
	return result, nil
}

/*
GetElementAtIndex retrieves the value at the specified index from the given slice.

Parameters:
- arr: The input slice from which the value needs to be retrieved.
- index: The index of the value to be fetched.

Returns:
- U: The value at the specified index in the slice. If the index is out of bounds, it returns the zero value of type U.

Example:

	arr := []int{1, 2, 3, 4, 5}
	val := GetElementAtIndex(arr, 2) // val will be 3

Note:
- The function returns the zero value of type U if the index is out of bounds.
*/
func GetElementAtIndex[U any, T []U](arr T, index int) U {
	var r U

	if index < 0 || index >= len(arr) {
		return r
	}

	r = arr[index]

	return r
}

// ArrayColumn extracts the values of a specified column (field) from a slice of structs.
// Generics must be type of slice that would be produced (ex: string, int)
func ArrayColumn[T any](slice interface{}, columnName string) ([]T, error) {
	// Validate that the provided slice is indeed a slice
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected a slice, but got %s", v.Kind())
	}

	var result []T

	// Iterate over the slice
	for i := 0; i < v.Len(); i++ {
		// Get the struct at the current index
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem() // Dereference if pointer
		}
		if item.Kind() != reflect.Struct {
			return nil, fmt.Errorf("expected struct type, but got %s", item.Kind())
		}

		// Get the field value by name
		field := item.FieldByName(columnName)
		if !field.IsValid() {
			return nil, fmt.Errorf("no such field: %s in struct", columnName)
		}

		// Try to cast the field value to the specified type
		fieldValue, ok := field.Interface().(T)
		if !ok {
			return nil, fmt.Errorf("field %s has incorrect type", columnName)
		}

		// Append the field value to the result
		result = append(result, fieldValue)
	}

	return result, nil
}

func ChunkArray[T any](arr []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

// Unique returns a new slice containing only the first occurrence
// of each element in input, preserving the original order.
//
// Type parameter T must be comparable so it can be used as a map key
// (e.g. string, int, bool, or structs with comparable fields).
//
// Time complexity: O(n)
// Space complexity: O(n)
//
// Example:
//
//	data := []string{"a", "b", "a", "c"}
//	result := Unique(data)
//	// result == []string{"a", "b", "c"}
func Unique[T comparable](input []T) []T {
	seen := make(map[T]struct{}, len(input))
	result := make([]T, 0, len(input))

	for _, v := range input {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// UniqueBy returns a new slice containing only the first occurrence
// of each element in input, determined by the value returned from keyFn.
//
// This is useful when T is not comparable (e.g. slices, maps)
// or when uniqueness should be based on a specific field.
//
// The order of elements is preserved based on their first appearance.
//
// Time complexity: O(n)
// Space complexity: O(n)
//
// Example:
//
//	type User struct {
//	    ID   int
//	    Name string
//	}
//
//	users := []User{
//	    {ID: 1, Name: "Alice"},
//	    {ID: 1, Name: "Bob"},
//	    {ID: 2, Name: "Charlie"},
//	}
//
//	result := UniqueBy(users, func(u User) int {
//	    return u.ID
//	})
//
//	// result == []User{
//	//   {ID: 1, Name: "Alice"},
//	//   {ID: 2, Name: "Charlie"},
//	// }
func UniqueBy[T any, K comparable](input []T, keyFn func(T) K) []T {
	seen := make(map[K]struct{})
	result := []T{}

	for _, v := range input {
		k := keyFn(v)
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

/**
 * AddSingleQuoteToEachString adds single quotes to each string in the input slice.
 * e.g. ["a", "b", "c"] becomes ["'a'", "'b'", "'c'"]
 */
func AddSingleQuoteToEachString(input []string) []string {
	result := make([]string, len(input))
	for i, str := range input {
		result[i] = "'" + str + "'"
	}

	return result
}
