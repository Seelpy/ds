package unique

type ReadOnlyTextCounter interface {
	GetCount(text string) (int, error)
}

type TextCounter interface {
	ReadOnlyTextCounter
	Inc(text string) error
	Dec(text string) error
}
