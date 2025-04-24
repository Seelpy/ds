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
	Get(textID uuid.UUID) (TextStatistics, error)
}

type TextStatisticsRepository interface {
	ReadOnlyTextStatisticsRepository
	Store(statistics TextStatistics) error
	Remove(textID uuid.UUID) error
}
