package handler

func ifNotEmpty(s string, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func ifNotNil[T any](v *T, fallback T) T {
	if v == nil {
		return fallback
	}
	return *v
}
