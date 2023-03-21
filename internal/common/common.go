package common

import "sort"

func SortMapByKeys[T any](items map[string]T) []T {
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var result []T
	for _, k := range keys {
		result = append(result, items[k])
	}
	return result
}
