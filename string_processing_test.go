package scholarScraper

import (
	"testing"
)

func TestRemoveForbiddenChars(test *testing.T) {
	testCases := map[string]string {
		"sample string"		: "sample string",
		"<sample string>"	: "sample string",
		"<is it sample string?>" : "is it sample string",
		"<is.it.sample.string?>" : "isitsamplestring",
		"\\<is.it.\\sample.\\string?>" : "isitsamplestring",
		"\"<**yeah. | It is, actually, a sample string!**>\"" : "yeah  It is actually a sample string!",
		"'/|;;\\\"*<:string, string..?:>*\"/|;;\\'" : "'string string'",
	}

	for str, expectedResult := range testCases {
		if actualResult := removeForbiddenChars(str); actualResult != expectedResult {
			return t.Error("Mismatch!\nExpected:", expectedResult, "\nGot:", actualResult)
		}
	}
}
