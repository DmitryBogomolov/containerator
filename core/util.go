package core

// TransformSlice transforms slice elements and makes new slice.
func TransformSlice[T any, P any](list []T, fn func(T) P) []P {
	if list == nil {
		return nil
	}
	result := make([]P, len(list))
	for i, item := range list {
		result[i] = fn(item)
	}
	return result
}
