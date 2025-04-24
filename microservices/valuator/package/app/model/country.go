package model

import (
	"errors"
	"github.com/gofrs/uuid"
)

type Country string

var (
	ErrCountryNotFound = errors.New("country not found")
)

const (
	RussiaCountry  Country = "Russia"
	FranceCountry  Country = "France"
	GermanyCountry Country = "Germany"
	UAECountry     Country = "UAE"
	IndiaCountry   Country = "India"
)

type TextCountryRepository interface {
	Store(textID uuid.UUID, country Country) error
}
