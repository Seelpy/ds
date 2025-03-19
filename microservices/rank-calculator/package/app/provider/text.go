package provider

import (
	"errors"
	"github.com/gofrs/uuid"
)

var (
	ErrTextNotFound = errors.New("text not found")
)

type TextData struct {
	ID    uuid.UUID
	Value string
}

type TextProvider interface {
	Get(id uuid.UUID) (TextData, error)
}
