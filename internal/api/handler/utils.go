package handler

func ifNotNil[T any](v *T, fallback T) T {
	if v == nil {
		return fallback
	}
	return *v
}
