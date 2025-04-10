package calculator

import (
	"rankcalculator/package/app/unique"
	"unicode/utf8"
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

type TextStatistics struct {
	AlphabetCount int
	AllCount      int
	IsDuplicate   bool
}

func NewRankCalculator(textCounter unique.ReadOnlyTextCounter) RankCalculator {
	return RankCalculator{
		textCounter: textCounter,
	}
}

type RankCalculator struct {
	textCounter unique.ReadOnlyTextCounter
}

func (c *RankCalculator) Calculate(text string) (TextStatistics, error) {
	count, err := c.textCounter.GetCount(text)
	if err != nil {
		return TextStatistics{}, err
	}

	alphabetSymbolsCount, allCount := c.symbolStatistics(text)

	return TextStatistics{
		AlphabetCount: alphabetSymbolsCount,
		AllCount:      allCount,
		IsDuplicate:   count > 0,
	}, nil
}

func (c *RankCalculator) symbolStatistics(text string) (alphabetSymbolsCount, allCount int) {
	for _, sym := range text {
		allCount++
		if alphabetMap[sym] {
			alphabetSymbolsCount++
		}
	}
	return
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
