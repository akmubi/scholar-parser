package scholarParser

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

// "@sample string/hello !!!\u32a7" -> "sample-string-hello-"
func MakeStringPretty(source string) (result string) {
	// source = removeForbiddenChars(source)
	for _, rune := range source {
		if (unicode.IsLetter(rune) || unicode.IsDigit(rune)) {
			result += string(rune)
		} else if unicode.IsSpace(rune) || rune == '\\' || rune == '/' {
			result += "-"
		} else {
			result += ""
		}
	}

	// " -sample-string-hello--" -> "sample-string-hello" 
	result = strings.Trim(result, " -")
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


// 4 -> { "0", "10", "30", "40"}
func makeIntStringSlice(upperBorder int64) (result []string) {
	for i := int64(0); i <= upperBorder; i++ {
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

// Что+такое+жизнь+-**".~

// allowed charactes: '-', '_', '*', '"', '~', '.'
func replaceNonLetters(source string) (result string) {

	for _, rune := range source {
		if	((rune == '-') || (rune == '*') || (rune == '_') ||
			(rune == '"') || (rune == '~') || (rune == '.')) ||
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