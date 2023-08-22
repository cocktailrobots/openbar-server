package functional

func Map[T, U any](f func(T) U, ts []T) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}

	return us
}
