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
			test.Errorf("Mismatch!\nExpected: '%s'\nGot: '%s'", expectedResult, actualResult)
		}
	}
}


func TestDomainName(test *testing.T) {
	testCases := map[string]string {
		"https://github.com/anaskhan96/soup" : "https://github.com/",
		"https://scholar.google.com/scholar?q=cyber-physical+system+design" : "https://scholar.google.com/",
		"https://vk.com/" : "https://vk.com/",
		"http://some.site/basic/knowledge/understood/" : "http://some.site/",
		"" : "",
		"someRandomString" : "someRandomString",
	}

	for url, extected := range testCases {
		if actual := domainName(url); actual != extected {
			test.Error("Domain name mismatch!\nExpected:", extected, "\nGot:", actual)
		}
	}
}