package main

// helpers for the parse

func toSlice[T any](v any) []T {
	if v == nil {
		return nil
	}

	// if it's already the right type, just return it
	if _, ok := v.([]T); ok {
		return v.([]T)
	}

	originalSlice := v.([]any)
	newSlice := make([]T, len(originalSlice))

	for i, val := range originalSlice {
		newSlice[i] = val.(T)
	}

	return newSlice
}

func toValue[T any](v any) T {
	if v == nil {
		var zero T
		return zero
	}

	return v.(T)
}
