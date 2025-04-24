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

func (s *StatisticsService) RankText(textID uuid.UUID) (provider.TextData, model.TextStatistics, error) {
	text, err := s.provider.Get(textID)
	if err != nil {
		return provider.TextData{}, model.TextStatistics{}, err
	}

	statistics, err := s.calculator.Calculate(text.Value)
	log.Printf("Afdaksoj%v", statistics)
	if err != nil {
		return provider.TextData{}, model.TextStatistics{}, err
	}
	err = s.repo.Store(model.TextStatistics{
		TextID:           textID,
		AllAlphabetCount: statistics.AlphabetCount,
		AllCount:         statistics.AllCount,
		IsDuplicate:      statistics.IsDuplicate,
	})
	if err != nil {
		return provider.TextData{}, model.TextStatistics{}, err
	}
	err = s.counter.Inc(text.Value)
	if err != nil {
		return provider.TextData{}, model.TextStatistics{}, err
	}
	stat, err := s.repo.Get(textID)
	return text, stat, err
}

func (s *StatisticsService) GetStatistics(textID uuid.UUID) (model.TextStatistics, error) {
	return s.repo.Get(textID)
}

func (s *StatisticsService) RemoveStatistics(textID uuid.UUID, textValue string) error {
	err := s.repo.Remove(textID)
	if err != nil {
		return err
	}

	return s.counter.Dec(textValue)
}
