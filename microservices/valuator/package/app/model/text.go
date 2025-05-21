package model

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/mono83/maybe"
	"valuator/package/app/authorization"
)

var (
	ErrTextNotFound         = errors.New("text not found")
	ErrTextPermissionDenied = errors.New("text permission denied")
)

type TextID uuid.UUID

type Text interface {
	ID() TextID
	UserID() uuid.UUID
	Value() string
}

type TextReadRepository interface {
	Find(ctx authorization.Context, id TextID) (maybe.Maybe[Text], error)
	ListAll() ([]Text, error)
}

type TextUpdateRepository interface {
	Store(ctx authorization.Context, text Text) error
	Create(ctx authorization.Context, value string) Text
	Remove(ctx authorization.Context, text Text) error
}

type TextRepository interface {
	TextUpdateRepository
	TextReadRepository
}

func LoadText(id TextID, userID uuid.UUID, value string) Text {
	return &text{
		id:     id,
		userID: userID,
		value:  value,
	}
}

func NewText(userID uuid.UUID, value string) Text {
	return &text{
		id:     TextID(uuid.NewV1()),
		userID: userID,
		value:  value,
	}
}

type text struct {
	id     TextID
	userID uuid.UUID
	value  string
}

func (t *text) ID() TextID {
	return t.id
}

func (t *text) UserID() uuid.UUID {
	return t.userID
}

func (t *text) Value() string {
	return t.value
}
