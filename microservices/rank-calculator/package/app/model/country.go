package model

import (
	"errors"
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
