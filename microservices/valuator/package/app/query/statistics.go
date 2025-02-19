package query

import (
	"github.com/gofrs/uuid"
	"unicode/utf8"
	"valuator/package/app/model"
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

func NewStatisticsQueryService(repo model.TextReadRepository) StatisticsQueryService {
	return &statisticsQueryService{
		repo: repo,
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
	repo model.TextReadRepository
}

func (s *statisticsQueryService) GetSummary(textID uuid.UUID) (TextStatistics, error) {
	text, err := s.repo.Find(model.TextID(textID))
	if err != nil {
		return TextStatistics{}, err
	}
	if text.IsEmpty() {
		return TextStatistics{}, model.ErrTextNotFound
	}
	all, err := s.repo.ListAll()
	if err != nil {
		return TextStatistics{}, err
	}

	return TextStatistics{
		SymbolStatistics: s.SymbolStatistics(text.Value()),
		UniqueStatistics: s.UniqueStatistic(text.Value(), all),
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

func (s *statisticsQueryService) UniqueStatistic(text model.Text, texts []model.Text) UniqueStatistics {
	for _, otherText := range texts {
		if otherText.ID() == text.ID() {
			continue
		}
		if otherText.Value() == text.Value() {
			return UniqueStatistics{
				IsDuplicate: true,
			}
		}
	}
	return UniqueStatistics{
		IsDuplicate: false,
	}
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
