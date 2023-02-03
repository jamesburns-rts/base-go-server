package ptr

func To[T any](t T) *T {
	return &t
}

func From[T any](t *T) (rt T) {
	if t != nil {
		return *t
	}
	return rt
}
