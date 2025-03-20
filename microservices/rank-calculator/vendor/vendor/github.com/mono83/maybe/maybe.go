package maybe

// Maybe is generic monadic container for value or it's absence
type Maybe[T any] struct {
	empty bool
	value T
}

// IsEmpty returns true if container has no data
func (m Maybe[T]) IsEmpty() bool { return m.empty }

// IsPresent returns true if container has data
func (m Maybe[T]) IsPresent() bool { return !m.empty }

// Value returns container value as-is. It is not recommended
// calling this method, use Maybe.Or, Maybe.OrElse, Maybe.OrElseGet
// or Maybe.Get instead.
func (m Maybe[T]) Value() T { return m.value }

// Get returns value and is present flag.
func (m Maybe[T]) Get() (T, bool) {
	return m.value, !m.empty
}

// Filter applies given predicate of container returning original
// value only if predicate resulted in true.
func (m Maybe[T]) Filter(predicate func(T) bool) Maybe[T] {
	if m.IsEmpty() || predicate == nil || !predicate(m.Value()) {
		return Nothing[T]()
	}
	return m
}

// IfPresent invokes given consumer callback only if container has value.
func (m Maybe[T]) IfPresent(consume func(T)) Maybe[T] {
	if m.IsPresent() && consume != nil {
		consume(m.Value())
	}
	return m
}

// IfPresent invokes given consumer callback if value present and plain
// func if container empty.
func (m Maybe[T]) IfPresentOrElse(consume func(T), els func()) Maybe[T] {
	if m.IsPresent() {
		if consume != nil {
			consume(m.Value())
		}
	} else {
		if els != nil {
			els()
		}
	}
	return m
}

// Or produces Maybe:
// 1. If current not empty - return it
// 2. If current is empty - return produced by supply func
// 3. Nothing if current is empty and supply function is nil
func (m Maybe[T]) Or(supply func() Maybe[T]) Maybe[T] {
	if m.IsPresent() {
		return m
	} else if supply == nil {
		return Nothing[T]()
	}
	return supply()
}

// OrElse returns value from Maybe if it presents, otherwise
// given value is returned.
func (m Maybe[T]) OrElse(other T) T {
	if m.IsEmpty() {
		return other
	}
	return m.Value()
}

// OrElseGet returns value from Maybe if it presents, otherwise
// value produced by supply function returnes.
func (m Maybe[T]) OrElseGet(supply func() T) T {
	if m.IsPresent() {
		return m.Value()
	}

	return supply()
}
