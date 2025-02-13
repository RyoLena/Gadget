package Gadget

func ToPtr[T any](v T) *T {
	return &v
}
