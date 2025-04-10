package maybe

import "reflect"

// Nilable constructs new Maybe and performs emptiness
// check of given data. Will produce empty Maybe for following cases:
// - nil value
// - zero-length slice
// - zero-length map
func Nilable[T any](value T) Maybe[T] {
	if isNotNil(value) {
		return Just(value)
	}

	return Nothing[T]()
}

// FilterNotNil filters Maybe leaving only non-nil value.
func (m Maybe[T]) FilterNotNil() Maybe[T] {
	return m.Filter(isNotNil[T])
}

func isNotNil[T any](value T) bool {
	v := reflect.ValueOf(value)
	k := v.Kind()
	if k == reflect.Invalid {
		return false
	}
	switch k {
	case reflect.Slice, reflect.Map:
		if v.IsNil() || v.Len() == 0 {
			return false
		}
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Pointer, reflect.UnsafePointer:
		if v.IsNil() {
			return false
		}
	}
	return true
}
