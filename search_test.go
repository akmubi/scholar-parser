package scholarScraper

import (
	"testing"
)

func TestPrefixTable(test *testing.T) {
	testCases := map[string][]int {
		"TEST"			: { 0, 0, 0, 1 },
		"AAAA"			: { 0, 1, 2, 3 },
		"AABA"			: { 0, 1, 0, 1 },
		"AAAAAAAAAAAB"	: { 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0 },
		"ABCABCABC"		: { 0, 0, 0, 1, 2, 3, 4, 5, 6 },
		"AAAABAAABA"	: { 0, 1, 2, 3, 0, 1, 2, 3, 0, 1 },
		"ABCDEF"		: { 0, 0, 0, 0, 0, 0 },
		"AABAACAABAA"	: { 0, 1, 0, 1, 2, 0, 1, 2, 3, 4, 5 },
		"AAACAAAAAC"	: { 0, 1, 2, 0, 1, 2, 3, 3, 3, 4 },
		"AAABAAA"		: { 0, 1, 2, 0, 1, 2, 3 },
	}

	if _, err := prefixTable(""); err == nil {
		test.Error("There was expected an error, but it hasn't appeared")
	}

	for pattern, tableExpected := range testCases {
		tableActual, _ := prefixTable(pattern)
		if len(tableActual) != len(tableExpected){
			test.Error("Wrong table length was taken")
		}

		for i := range tableActual {
			if tableActual[i] != tableExpected[i] {
				test.Error("Mismatch!\nExpected:", tableExpected, "\nGot:", tableActual)
				break
			}
		}
	}
}

func TestFindPattern(test *testing.T) {

	if _, err := FindAbstractIndex("Some text", ""); err == nil {
		test.Error("Pattern is empty. There was expected an error, but it hasn't appeared")
	}

	if _, err := FindAbstractIndex("", "Some pattern"); err == nil {
		test.Error("Text is empty. There was expected an error, but it hasn't appeared")
	}

	if _, err := FindAbstractIndex("", ""); err == nil {
		test.Error("Both text and pattern are empty. There was expected an error, but it hasn't appeared")
	}

	if index, _ := FindAbstractIndex("\nAbstract-", "Abstract"); index != 1 {
		test.Error("Mismatch!\nExpected:", 1, "\nGot:", index)
	}
	if index, _ := FindAbstractIndex("\nAbstract|", "\nAbstract"); index != 0 {
		test.Error("Mismatch!\nExpected:", 0, "\nGot:", index)
	}

	if index, _ := FindAbstractIndex("\nAbstracting", "\nAbstract"); index != -1 {
		test.Error("Mismatch!\nExpected:", -1, "\nGot:", index)
	}

	if index, _ := FindAbstractIndex("\nAbstraction", "Abstract"); index != -1 {
		test.Error("Mismatch!\nExpected:", -1, "\nGot:", index)
	}

	if index, _ := FindAbstractIndex("\nAbstrac", "Abstract"); index != -1 {
		test.Error("Mismatch!\nExpected:", -1, "\nGot:", index)
	}

	if index, _ := FindPatternIndex("\nXII. Conclusion", "Conclusion"); index != 6 {
		test.Error("Mismatch!\nExpected:", 6, "\nGot:", index)
	}

	if index, _ := FindPatternIndex("\nXII. Conclusions And Future Work", "Conclusion"); index != 6 {
		test.Error("Mismatch!\nExpected:", 6, "\nGot:", index)
	}
}