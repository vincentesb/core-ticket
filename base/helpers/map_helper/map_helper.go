package map_helper

// Map sued to create a map from an array
func Map[T any, V comparable](src []T, key func(T) V) map[V]T {
	var result = make(map[V]T)
	for _, v := range src {
		result[key(v)] = v
	}
	return result
}
