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

type User struct {
	UserID       uuid.UUID
	Login        string
	PasswordHash string
	Country      Country
}

type UserRegisterInput struct {
	Login        string
	PasswordHash string
	Country      Country
}

type UserRepository interface {
	Register(user UserRegisterInput) (User, error)
	GetByID(userID uuid.UUID) (User, error)
	GetByLogin(login string) (User, error)
}
