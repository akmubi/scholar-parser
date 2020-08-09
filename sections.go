package scholarScraper

import (
	"io/ioutil"
	"errors"
	"fmt"
	"strings"
)

type AbstractSection struct {
	filepath				string
	startIndex				int
	keywordsIndex			int
	introductionIndex		int
}

func (s *AbstractSection) New() {
	s.filepath				= ""
	s.startIndex			= -1
	s.keywordsIndex			= -1
	s.introductionIndex		= -1
}

var (
	// Abstract
	abstractSamples = [...]string{
		"\nAbstract",
		"\nABSTRACT",
		"\nA BSTRACT",
		"\nA B S T R A C T",
	}

	// Keywords
	keywordsSamples = [...]string{
		"\nKeywords",
		"\nKEYWORDS",
		"\nK EYWORDS",
		"\nK E Y W O R D S",
		"\nIndex Terms",
		"\nINDEX TERMS",
	}

	// Introduction
	introductionSamples = [...]string{
		"Introduction",
		"INTRODUCTION",
		"I NTRODUCTION",
		"I N T R O D U C T I O N",
	}
)

func (s AbstractSection) Print() {
	fmt.Println("Start:", s.startIndex, "| Keywords:", s.keywordsIndex, "| Introduction:", s.introductionIndex)
}

func (s *AbstractSection) GetContent() (string, error) {

	if s.startIndex == -1 || (s.keywordsIndex == -1 && s.introductionIndex == -1) {
		return "", errors.New("No Abstract")
	}

	txtFile, err := ioutil.ReadFile(s.filepath)
	if err != nil {
		return "", err
	}

	if len(txtFile) == 0 {
		return "", errors.New("Empty file")
	}

	text := string(txtFile)

	var kwrdSection, intrdSection string

	// find abstract ~ keywords section
	if s.keywordsIndex != -1 {
		for i, char := range text {
			if i >= s.startIndex && i <= s.keywordsIndex {
				kwrdSection += string(char)
			}
		}
		kwrdSection = strings.TrimSpace(kwrdSection)

		if kwrdSection != "" {
			// put period if it's a letter
			if isLetter(kwrdSection[len(kwrdSection) - 1]) {
				kwrdSection += "."
			}
			fmt.Println("[Easy. Keywords]")
			return kwrdSection, nil
		}

		// keywordsIndex is not valid anymore
		s.keywordsIndex = -1
		fmt.Println("***Keywords too close to Abstract***")
	}

	// if there were no letters between abstract and keywords or 
	// if there were no keywords at all
	// then find abstract ~ introduction 
	if s.introductionIndex != -1 {
		for i, char := range text {
			if i >= s.startIndex && i <= s.introductionIndex {
				intrdSection += string(char)
			}
		}
		intrdSection = strings.TrimSpace(intrdSection)
		intrdSection = strings.TrimSuffix(intrdSection, "1.")
		intrdSection = strings.TrimSuffix(intrdSection, "I.")

		if intrdSection != "" {
			fmt.Println("[Yeeah. Introduction]")
			return intrdSection, nil
		}
		// introductionIndex is not valid anymore too
		s.introductionIndex = -1
		fmt.Println("***Introdution too close to Abstract***")
	}

	return "", errors.New("----What? Looks like keywords were too close to 'Abstract' and introduction wasn't found or wise versa----")
}

func readFile(filepath string) (text string, err error) {
	txtFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	text = string(txtFile)
	if len(text) == 0 {
		return "", errors.New("Empty file")
	}
	return text, nil
}

func FindAbstract(filepath string) (abstract AbstractSection, err error) {
	abstract.New()

	text, err := readFile(filepath)
	if err != nil {
		return abstract, err
	}

	abstract.filepath = filepath

	// find 'Abstract' word
	var abstrIndex int
	for _, abstr := range abstractSamples {
		if abstrIndex, _ = FindAbstractIndex(text, abstr); abstrIndex != -1 {
			// skip 'Abstract' word
			abstrIndex += len(abstr)
			abstract.startIndex = abstrIndex
			break
		}
	}

	if abstract.startIndex != -1 {
		// abstract is found. now try to find keywords
		var keywrdIndex int
		for _, keywrd := range keywordsSamples {
			if keywrdIndex, _ = FindPatternIndex(text, keywrd); keywrdIndex != -1 {
				abstract.keywordsIndex = keywrdIndex
				break
			}
		}

		// and introduction just in case
		var intrdIndex int
		for _, intrd := range introductionSamples {
			if intrdIndex, _ = FindPatternIndex(text, intrd); intrdIndex != -1 {
				// "Abstract: ... I" -> "Abstract: ... "
				abstract.introductionIndex = intrdIndex - 1
				break
			}
		}

		if abstract.keywordsIndex == -1 && abstract.introductionIndex == -1 {
			return AbstractSection{ filepath, -1, -1, -1 }, errors.New("Coundn't find abstract ending")
		}
	} else {
		return AbstractSection{ filepath, -1, -1, -1 }, errors.New("Coundn't find abstract")
	}
	return abstract, nil
}

func SplitSection(rawAbstract string) string {
	// ". | Begining of abstract..." ->
	// "Begining of abstract..."
	rawAbstract = strings.TrimLeftFunc(rawAbstract, func(r rune)bool {
		return !(r >= 'A' && r <= 'Z') && !(r >= 'a' && r <= 'z')
	})

	// "...ending of abstract. " ->
	// "...ending of abstract."
	rawAbstract = strings.TrimRightFunc(rawAbstract, func(r rune)bool {
		return !(r >= 'A' && r <= 'Z') && !(r >= 'a' && r <= 'z') && r != '.'
	})

	// replace all new lines with space
	var result string
	lines := strings.Split(rawAbstract, "\n")
	for i, line := range lines {
		result += strings.TrimSpace(line)
		if i != len(lines) - 1 {
			result += " "
		}
	}
	
	// and remove double spaces	
	result = strings.Replace(result, "  ", " ", -1)
	return result
} 