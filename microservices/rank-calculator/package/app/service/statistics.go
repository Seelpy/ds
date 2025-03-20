package service

import (
	"github.com/gofrs/uuid"
	"rankcalculator/package/app/calculator"
	"rankcalculator/package/app/model"
)

func NewStatisticsService(repo model.TextStatisticsRepository, calculator calculator.RankCalculator) StatisticsService {
	return StatisticsService{
		repo:       repo,
		calculator: calculator,
	}
}

type StatisticsService struct {
	repo       model.TextStatisticsRepository
	calculator calculator.RankCalculator
}

func (s *StatisticsService) RankText(textID uuid.UUID) error {
	statistics, err := s.calculator.Calculate(textID)
	if err != nil {
		return err
	}
	return s.repo.Store(model.TextStatistics{
		TextID:           textID,
		AllAlphabetCount: statistics.AlphabetCount,
		AllCount:         statistics.AllCount,
		IsDuplicate:      statistics.IsDuplicate,
	})
}

func (s *StatisticsService) RemoveStatistics(textID uuid.UUID) error {
	return s.repo.Remove(textID)
}
