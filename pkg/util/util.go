package util

func Ptr[T any](in T) *T {
	return &in
}
