package unique

type TextUniquenessChecker interface {
	IsUnique(text string) (bool, error)
}

type TextUniquenessStorage interface {
	TextUniquenessChecker
	Store(text string) error
}
