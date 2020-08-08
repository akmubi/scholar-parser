package scholarScraper

import (
	"errors"
)

func defaultCheckLastChar(char byte) bool {
	return false
}

func isLetter(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func FindAbstractIndex(text, pattern string) (index int, err error) {
	return findPattern(text, pattern, isLetter)
}

func FindPatternIndex(text, pattern string) (index int, err error) {
	return findPattern(text, pattern, defaultCheckLastChar)
}

func findPattern(text, pattern string, checkLastChar func(char byte) bool) (index int, err error) {
	
	textLength := len(text)
	if textLength == 0 {
		return -1, errors.New("Text is empty")
	}

	patternLength := len(pattern)
	
	table, err := prefixTable(pattern)
	if err != nil {
		return -1, err
	}

	var i, j int
	for i < textLength {
		if text[i] == pattern[j] {
			i++
			j++
		}
		if j == patternLength {
			if i < textLength && checkLastChar(text[i]) {
				// Example: "\nAbstraction" (letter 'i' after "\nAbstract")
				// revert j and start searching after 'i'
				i++
				j = 0
				continue
			}
			return i - j, nil
		} else if i < textLength && pattern[j] != text[i] {
			if j != 0 {
				j = table[j - 1]
			} else {
				i++
			}
		}
	}
	return -1, nil
}

func prefixTable(pattern string) (table []int, err error) {
	patternLength := len(pattern)
	if patternLength == 0 {
		return nil, errors.New("Pattern is Empty");
	}

	table = make([]int, patternLength)

	i, j := 0, 1
	table[0] = 0

	for j < patternLength {
		if pattern[i] == pattern[j] {
			table[j] = i + 1
			i++
			j++
		} else if i == 0 {
			table[j] = 0
			j++
		} else {
			i = table[i - 1]
		}
	}
	return table, nil
}