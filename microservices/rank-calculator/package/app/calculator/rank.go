package calculator

import (
	"github.com/gofrs/uuid"
	"rankcalculator/package/app/provider"
)

type TextStatistics struct {
	TextID        uuid.UUID
	AlphabetCount int
	AllCount      int
	IsDuplicate   bool
}

func NewRankCalculator(textProvider provider.TextProvider) RankCalculator {
	return RankCalculator{
		textProvider: textProvider,
	}
}

type RankCalculator struct {
	textProvider provider.TextProvider
}

func (c *RankCalculator) Calculate(id uuid.UUID) (TextStatistics, error) {
	_, err := c.textProvider.Get(id)
	if err != nil {
		return TextStatistics{}, err
	}

	return TextStatistics{
		TextID:        uuid.UUID{},
		AlphabetCount: 0,
		AllCount:      0,
		IsDuplicate:   false,
	}, nil
}
