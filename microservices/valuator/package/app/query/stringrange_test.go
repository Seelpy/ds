package query

import (
	"testing"
	"unicode/utf8"
)

func TestRu(t *testing.T) {
	str := "ПРИВЕТ"
	assertLen(t, getUtfRunes(str), getBasicRunes(str))
}

func TestEn(t *testing.T) {
	str := "HELLO"
	assertLen(t, getUtfRunes(str), getBasicRunes(str))
}

func TestNum(t *testing.T) {
	str := "123"
	assertLen(t, getUtfRunes(str), getBasicRunes(str))
}

func TestChina(t *testing.T) {
	str := "熊猫"
	assertLen(t, getUtfRunes(str), getBasicRunes(str))
}

func TestSmile(t *testing.T) {
	str := "☺㋡"
	assertLen(t, getUtfRunes(str), getBasicRunes(str))
}

func getUtfRunes(str string) []rune {
	utfRunes := make([]rune, 0)
	tmp := str
	for len(tmp) > 0 {
		r, size := utf8.DecodeRuneInString(tmp)
		tmp = tmp[size:]
		utfRunes = append(utfRunes, r)
	}

	return utfRunes
}

func getBasicRunes(str string) []rune {
	stdRunes := make([]rune, 0)
	for _, runeValue := range str {
		stdRunes = append(stdRunes, runeValue)
	}
	return stdRunes
}

func assertLen(t *testing.T, str1, str2 []rune) {
	if len(str1) != len(str2) {
		t.Errorf("%d != %d", len(str1), len(str2))
	}
}
