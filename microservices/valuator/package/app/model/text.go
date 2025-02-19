package model

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
)

var (
	ErrTextNotFound = errors.New("text not found")
)

type TextID uuid.UUID

type Text interface {
	ID() TextID
	Value() string
}

type TextReadRepository interface {
	Find(id TextID) (maybe.Maybe[Text], error)
	ListAll() ([]Text, error)
}

type TextUpdateRepository interface {
	Store(text Text) error
	Create(value string) Text
	Remove(text Text) error
}

type TextRepository interface {
	TextUpdateRepository
	TextReadRepository
}

func LoadText(id TextID, value string) Text {
	return &text{
		id:    id,
		value: value,
	}
}

func NewText(value string) Text {
	return &text{
		id:    TextID(uuid.NewV1()),
		value: value,
	}
}

type text struct {
	id    TextID
	value string
}

func (t *text) ID() TextID {
	return t.id
}

func (t *text) Value() string {
	return t.value
}
