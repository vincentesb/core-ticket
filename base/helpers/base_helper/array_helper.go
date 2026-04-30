package base_helper

func GetUniqueValues[T int | string | float64 | float32](source []T) []T {
	uniqueValues := make([]T, 0)
	uniqueMap := make(map[T]bool)
	for _, value := range source {
		if _, ok := uniqueMap[value]; !ok {
			uniqueMap[value] = true
			uniqueValues = append(uniqueValues, value)
		}
	}
	return uniqueValues
}
