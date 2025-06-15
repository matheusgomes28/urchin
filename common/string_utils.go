/*
    This file should provide utility functions for
	handling string-related containers.

	The context is to eventually support all the nice
	to haves that other languages have like filter, map,
	filter_map, etc.
*/

package common

func FilterStrings(strings []string, predicate func(string) bool) []string {
	result := make([]string, 0)
	for _, value := range strings {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

// This makes this file no longer string utils only!
func Map[K any, V any](inputs []K, mapper func(K) V) []V {
	result := make([]V, 0)
	for _, value := range inputs {
		result = append(result, mapper(value))
	}
	return result
}
