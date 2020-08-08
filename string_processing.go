package scholarScraper

import (
	"unicode"
	"strings"
	"strconv"
	"errors"
	"fmt"
)

func removeForbiddenChars(source string) (result string) {
	forbiddedChars := "<>:\"/\\|*?.,;"
	result = source
	for _, rune := range forbiddedChars {
		result = strings.ReplaceAll(result, string(rune), "")
	}
	return
}

// "@sample string/hello !!!\u32a7" -> "sample_string_hello"
func MakeStringPretty(source string) (result string) {
	// source = removeForbiddenChars(source)
	for _, rune := range source {
		if (unicode.IsLetter(rune) || unicode.IsDigit(rune)) {
			result += string(rune)
		} else if unicode.IsSpace(rune) || rune == '\\' || rune == '/' {
			result += "_"
		} else {
			result += ""
		}
	}

	// " _sample_string_hello__" -> "sample_string_hello" 
	result = strings.Trim(result, " _")
	return
}

// only for checking if given string is allowed language or not
func isAllowedLanguage(allowedLanguages string, language string) bool {
	
	splitted := strings.Split(allowedLanguages, "&")
	for _, allowedLanguage := range splitted {
		if allowedLanguage == language {
			return true
		}
	}
	return false
}


// 4 pages -> { "0", "10", "20", "30" }
func makeIntStringSlice(upperBorder int64) (result []string) {
	for i := int64(0); i < upperBorder; i++ {
		result = append(result, strconv.FormatInt(i * 10, 10))
	}
	return
}

func getFolder(filepath string) string {
	var i int
	var length = len(filepath)
	if (filepath[length - 1] == '/') || (filepath[length - 1] == '\\') {
		return filepath
	}
	for i = length - 1; i >= 0; i-- {
		if (filepath[i] == '/') || (filepath[i] == '\\') {
			return filepath[:i + 1]
		}
	}
	return ""
}

// allowed charactes: '-', '_', '*', '"', '~', '.', '<', '>'
func replaceNonLetters(source string) (result string) {

	for _, rune := range source {
		if	((rune == '-') || (rune == '*') || (rune == '_') ||
			(rune == '"') || (rune == '~') || (rune == '.')) ||
			(rune == '<') || (rune == '>') ||
			(unicode.IsLetter(rune)) {
			result += string(rune)
		} else {
			result += fmt.Sprintf("%%%X", int(rune))
		}
	}
	return
}

// Assume that source string has \ufffd, \uffff, ... letters
func unicodeConvert(source string) (result string) {
	for _, rune := range source {
		result += string(rune)
	}
	return
}

// searching author names in string
// example: (in search result)	C Cifuentes -> 
//			(in document)		Cristina Cifuentes
func FindAuthorName(pdfText string, authorInitials string) (fullname string, err error) {
	nameParts := strings.Split(authorInitials, " ")
	// initials = { "C" }
	// name = { "Cifuentes" }
	var name []string
	var initials string

	for i := range nameParts {
		var isInitials bool = true
		for _, letter := range nameParts[i] {
			if unicode.IsLower(letter) {
				isInitials = false
			}
		}
		if isInitials {
			initials = reverseString(nameParts[i])
		} else {
			name = append(name, nameParts[i])
		}
	}

	if len(name) == 0 || len(initials) == 0 {
		return "", errors.New("There is no name")
	}

	// initial = { "JC" },
	// name = { "Jensen" }
	nameIndex, err := FindPatternIndex(pdfText, name[0]);
	if err != nil {
		return "", err
	}
	if nameIndex == -1 {
		return "", errors.New("Couldn't find name in PDF file")
	}

	var startIndex int
	if nameIndex - 100 < 0 {
		startIndex = 0
	} else {
		startIndex = nameIndex - 100
	}

	// "Jeff C. Jensen ..." -> " .C ffeJ" 
	if nameIndex - 1 < 0 {
		return "", errors.New("Name without initials")
	}
	reversedTextSection := reverseString(pdfText[startIndex:nameIndex - 1])

	// " .C ffeJ"
	textRunes := []rune(reversedTextSection)
	
	// "CJ"
	initialRunes := []rune(initials)

	found := make([]string, len(initialRunes))
	var skipLetters int = 0
	for i := 0; i < len(initialRunes); i++ {
		for j := skipLetters; j < len(textRunes); j++ {
			if initialRunes[i] == textRunes[j] {
				for k := j; k >= 0 && unicode.IsLetter(textRunes[k]); k-- {
					found[i] += string(textRunes[k])
				}
				// before: " .C ffeJ", skipLetters = 0
				// after:  " ffeJ", skipLetters = 3
				skipLetters = j + 1
				// found. next
				break
			}
		}
	}

	for i := len(found) - 1; i >= 0; i-- {
		fullname += found[i] + " "
	}
	for i := range name {
		fullname += name[i] + " "
	}

	return strings.TrimRight(fullname, " "), nil
}

func reverseString(sample string) string {
	runes := []rune(sample)
	for i := 0; i < len(runes) / 2; i++ {
		temp := runes[i]
		runes[i] = runes[len(runes) - 1 - i]
		runes[len(runes) - 1 - i] = temp 
	}
	return string(runes)
}