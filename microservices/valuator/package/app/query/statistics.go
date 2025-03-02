package query

import (
	"github.com/gofrs/uuid"
	"unicode/utf8"
	"valuator/package/app/model"
	"valuator/package/app/unique"
)

const (
	lowerCaseEnAlphabet = "abcdefghijklmnopqrstuvwxyz"
	upperCaseEnAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerCaseRuAlphabet = "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	upperCaseRuAlphabet = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
	alphabet            = lowerCaseEnAlphabet + upperCaseEnAlphabet + lowerCaseRuAlphabet + upperCaseRuAlphabet
)

var (
	alphabetMap = generateAlphabetMap()
)

type StatisticsQueryService interface {
	GetSummary(textID uuid.UUID) (TextStatistics, error)
}

func NewStatisticsQueryService(repo model.TextReadRepository, counter unique.ReadOnlyTextCounter) StatisticsQueryService {
	return &statisticsQueryService{
		repo:    repo,
		counter: counter,
	}
}

type TextStatistics struct {
	SymbolStatistics
	UniqueStatistics
}

type SymbolStatistics struct {
	AlphabetCount int
	AllCount      int
}

type UniqueStatistics struct {
	IsDuplicate bool
}

type statisticsQueryService struct {
	repo    model.TextReadRepository
	counter unique.ReadOnlyTextCounter
}

func (s *statisticsQueryService) GetSummary(textID uuid.UUID) (TextStatistics, error) {
	text, err := s.repo.Find(model.TextID(textID))
	if err != nil {
		return TextStatistics{}, err
	}
	if text.IsEmpty() {
		return TextStatistics{}, model.ErrTextNotFound
	}
	uniqueStat, err := s.UniqueStatistic(text.Value())
	if err != nil {
		return TextStatistics{}, err
	}

	return TextStatistics{
		SymbolStatistics: s.SymbolStatistics(text.Value()),
		UniqueStatistics: uniqueStat,
	}, nil
}

func (s *statisticsQueryService) SymbolStatistics(text model.Text) SymbolStatistics {
	tmp := text.Value()
	result := make(map[rune]bool)
	alphabetSymbolsCount := 0
	allCount := 0
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
		result[r] = true
		allCount++
		if alphabetMap[r] {
			alphabetSymbolsCount++
		}
	}
	return SymbolStatistics{
		AlphabetCount: alphabetSymbolsCount,
		AllCount:      allCount,
	}
}

func (s *statisticsQueryService) UniqueStatistic(text model.Text) (UniqueStatistics, error) {
	count, err := s.counter.GetCount(text.Value())
	if err != nil {
		return UniqueStatistics{}, err
	}
	return UniqueStatistics{
		IsDuplicate: count > 1,
	}, nil
}

func generateAlphabetMap() map[rune]bool {
	result := make(map[rune]bool)
	tmp := alphabet
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
		result[r] = true
	}
	return result
}
