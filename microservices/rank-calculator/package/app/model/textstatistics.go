package model

import (
	"errors"
	"github.com/gofrs/uuid"
)

var (
	ErrStatisticsNotFound = errors.New("statistics not found")
)

type TextStatistics struct {
	TextID           uuid.UUID
	AllAlphabetCount int
	AllCount         int
	IsDuplicate      bool
}

func (t *TextStatistics) Rank() float64 {
	return 1 - (float64(t.AllAlphabetCount) / float64(t.AllCount))
}

type ReadOnlyTextStatisticsRepository interface {
	Get(userID uuid.UUID, textID uuid.UUID) (TextStatistics, error)
}

type TextStatisticsRepository interface {
	ReadOnlyTextStatisticsRepository
	Store(userID uuid.UUID, statistics TextStatistics) error
	Remove(userID uuid.UUID, textID uuid.UUID) error
}
