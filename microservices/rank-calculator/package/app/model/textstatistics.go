package model

import (
	"github.com/gofrs/uuid"
)

type TextStatistics struct {
	TextID           uuid.UUID
	AllAlphabetCount int
	AllCount         int
	IsDuplicate      bool
}

type TextStatisticsRepository interface {
	Store(statistics TextStatistics) error
	Remove(textID uuid.UUID) error
}
