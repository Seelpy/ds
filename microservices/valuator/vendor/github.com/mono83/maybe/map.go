package maybe

// Map constructs new Maybe container built using
// mapper function.
func Map[T any, R any](from Maybe[T], mapper func(T) R) Maybe[R] {
	if from.IsEmpty() || mapper == nil {
		return Nothing[R]()
	}

	return Nilable(mapper(from.Value()))
}

// nMap constructs new Maybe container built using
// mapper function but for non-nilable types only.
func nMap[T, R any](from Maybe[T], mapper func(T) R) Maybe[R] {
	if from.IsEmpty() || mapper == nil {
		return Nothing[R]()
	}

	return Just(mapper(from.Value()))
}

func (m Maybe[T]) MapToAny(mapper func(T) any) Maybe[any]             { return Map(m, mapper) }
func (m Maybe[T]) MapToString(mapper func(T) string) Maybe[string]    { return nMap(m, mapper) }
func (m Maybe[T]) MapToInt(mapper func(T) int) Maybe[int]             { return nMap(m, mapper) }
func (m Maybe[T]) MapToInt16(mapper func(T) int16) Maybe[int16]       { return nMap(m, mapper) }
func (m Maybe[T]) MapToInt32(mapper func(T) int32) Maybe[int32]       { return nMap(m, mapper) }
func (m Maybe[T]) MapToInt64(mapper func(T) int64) Maybe[int64]       { return nMap(m, mapper) }
func (m Maybe[T]) MapToUint(mapper func(T) uint) Maybe[uint]          { return nMap(m, mapper) }
func (m Maybe[T]) MapToUint16(mapper func(T) uint16) Maybe[uint16]    { return nMap(m, mapper) }
func (m Maybe[T]) MapToUint32(mapper func(T) uint32) Maybe[uint32]    { return nMap(m, mapper) }
func (m Maybe[T]) MapToUint64(mapper func(T) uint64) Maybe[uint64]    { return nMap(m, mapper) }
func (m Maybe[T]) MapToFloat32(mapper func(T) float32) Maybe[float32] { return nMap(m, mapper) }
func (m Maybe[T]) MapToFloat64(mapper func(T) float64) Maybe[float64] { return nMap(m, mapper) }
