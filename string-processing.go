package scholarScraper

import (
	"unicode"
	"strings"
	"strconv"
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
