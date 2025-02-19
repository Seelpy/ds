package model

type Text interface {
	ID() uui
}
type TextRepository interface {
	Store(text Text) error
}
