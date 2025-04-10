package service

import (
	"github.com/gofrs/uuid"
	"log"
	"rankcalculator/package/app/calculator"
	"rankcalculator/package/app/model"
	"rankcalculator/package/app/provider"
	"rankcalculator/package/app/unique"
)

func NewStatisticsService(repo model.TextStatisticsRepository, calculator calculator.RankCalculator, counter unique.TextCounter, provider provider.TextProvider) StatisticsService {
	return StatisticsService{
		repo:       repo,
		calculator: calculator,
		counter:    counter,
		provider:   provider,
	}
}

type StatisticsService struct {
	repo       model.TextStatisticsRepository
	calculator calculator.RankCalculator
	counter    unique.TextCounter
	provider   provider.TextProvider
}

func (s *StatisticsService) RankText(textID uuid.UUID) error {
	text, err := s.provider.Get(textID)
	if err != nil {
		return err
	}

	statistics, err := s.calculator.Calculate(text.Value)
	log.Printf("Afdaksoj%v", statistics)
	if err != nil {
		return err
	}
	err = s.repo.Store(model.TextStatistics{
		TextID:           textID,
		AllAlphabetCount: statistics.AlphabetCount,
		AllCount:         statistics.AllCount,
		IsDuplicate:      statistics.IsDuplicate,
	})
	if err != nil {
		return err
	}
	return s.counter.Inc(text.Value)
}

func (s *StatisticsService) RemoveStatistics(textID uuid.UUID, textValue string) error {
	err := s.repo.Remove(textID)
	if err != nil {
		return err
	}

	return s.counter.Inc(textValue)
}
