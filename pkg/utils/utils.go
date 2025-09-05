package utils

// Returns true if all elements in the slice satisfy the predicate.
func All[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if !pred(t) {
			return false
		}
	}
	return true
}

// Returns true if at least one element in the slice satisfies the predicate.
func Some[T any](ts []T, pred func(T) bool) bool {
	for _, t := range ts {
		if pred(t) {
			return true
		}
	}
	return false
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func FlatMap[A, B any](input []A, f func(A) []B) []B {
	var result []B
	for _, v := range input {
		result = append(result, f(v)...)
	}
	return result
}
