package maybe

// Nothing constructs new Maybe with no content.
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{empty: true}
}

// Just constructs new Maybe with given value.
// This constructor will create non-empty Maybe
// even for nil value, for protection use
// Nilable instead.
func Just[T any](value T) Maybe[T] {
	return Maybe[T]{value: value}
}

// Ptr constructs Maybe from pointer value.
// If nil given returns Nothing.
func Ptr[T any](value *T) Maybe[T] {
	if value == nil {
		return Nothing[T]()
	}
	return Just(*value)
}
