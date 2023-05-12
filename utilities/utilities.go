package utilities

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](items ...T) T {
	var result T
	if len(items) == 0 {
		return result
	}

	result = items[0]
	for _, item := range items {
		if item > result {
			result = item
		}
	}
	return result
}
