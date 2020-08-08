package scholarScraper

import (
	"testing"
)

func TestSplitInfo(test *testing.T) {
	testCases := map[string]info {
		"P Derler, EA Lee, S Tripakis, M Törngren - … on Cyber-Physical Systems, 2013 - dl.acm.org" : {
			authors : []string { "P Derler", "EA Lee", "S Tripakis", "M Törngren" },
			journalName : "on Cyber-Physical Systems",
			website : "dl.acm.org",
			isJournalNameLonger : true,
			year : 2013,
		},
		"S Zeadally, N Jabeur - 2016 - dl.acm.org" : {
			authors : []string { "S Zeadally", "N Jabeur" },
			journalName : "",
			website : "dl.acm.org",
			isJournalNameLonger : false,
			year : 2016,
		},
		"P Nuzzo, M Lora, YA Feldman… - 2018 Design …, 2018 - ieeexplore.ieee.org" : {
			authors : []string { "P Nuzzo", "M Lora", "YA Feldman" },
			journalName : "2018 Design",
			website : "ieeexplore.ieee.org",
			isJournalNameLonger : true,
			year : 2018,
		},
		"L Pike - 2016 IEEE Cybersecurity Development (SecDev), 2016 - ieeexplore.ieee.org" : {
			authors : []string { "L Pike" },
			journalName : "2016 IEEE Cybersecurity Development (SecDev)",
			website : "ieeexplore.ieee.org",
			isJournalNameLonger : false,
			year : 2016,
		},
		"W Wolf - Computer, 2009 - computer.org" : {
			authors : []string { "W Wolf" },
			journalName : "Computer",
			website : "computer.org",
			isJournalNameLonger : false,
			year : 2009,
		},
		"KE Mary Reena, A Theckethil Mathew… - … Problems in Engineering, 2015 - hindawi.com" : {
			authors : []string { "KE Mary Reena", "A Theckethil Mathew" },
			journalName : "Problems in Engineering",
			website : "hindawi.com",
			isJournalNameLonger : true,
			year : 2015,
		},
	}

	for info, expectedResult := range testCases {
		var article Article
		article.Info = info
		actualResult, err := article.splitInfo()
		if err != nil {
			test.Error("An error occured", err)
		}
		if len(expectedResult.authors) != len(actualResult.authors) {
			test.Error("Authors aren't match\nExpected:\n", expectedResult.authors, "\nGot:\n", actualResult.authors)
		}
		if expectedResult.journalName != actualResult.journalName {
			test.Error("Journal Names aren't match\nExpected:", expectedResult.journalName, "\nGot:", actualResult.journalName)
		}
		if expectedResult.website != actualResult.website {
			test.Error("Websites aren't match\nExpected:", expectedResult.website, "\nGot:", actualResult.website)
		}
		if expectedResult.isJournalNameLonger != actualResult.isJournalNameLonger {
			test.Error("isJournalNameLonger fields aren't match\nExpected:", expectedResult.isJournalNameLonger, "\nGot:", actualResult.isJournalNameLonger)
		}
		if expectedResult.year != actualResult.year {
			test.Error("Years aren't match\nExpected:", expectedResult.year, "\nGot:", actualResult.year)
		}
	}
}

func TestFindAuthorName(test *testing.T) {
	testCases := map[ [2]string ]string {
		{ "Jeff C. Jensen ; Danica H. Chang", "JC Jensen" } : "Jeff C Jensen",
		{ "Siddharth Sridhar ; Adam Hahn ; Manimaran Govindarasu", "M Govindarasu" } : "Manimaran Govindarasu",
		{ "Xiping Hu ; Terry H. S. Chu ; Henry C. B. Chan", "THS Chu" } : "Terry H S Chu",
		{ "Xiping Hu ; Terry H. S. Chu ; Henry C. B. Chan", "HCB Chan" } : "Henry C B Chan",
		{ "Like Mike; Rob. Pike", "LMR Pike" } : "Like Mike Rob Pike",
	}

	for args, expectedName := range testCases {
		actualName, err := FindAuthorName(args[0], args[1])
		if err != nil {
			test.Error("An error occured", err)
		}
		if actualName != expectedName {
			test.Errorf("Name mismatch\nExpected: '%s'\nGot: '%s'\n", expectedName, actualName)
		}
	}
}

// func TestMain(test *testing.T) {

// 	// Your query
// 	query := "Reverse compilation techniques"

// 	// configuration
// 	var config Config

// 	// set searching params
// 	parameters := map[string]string {
// 		"query"	: query,
// 		"lang"	: "en",
// 		"pages" : "1",
// 	}
	
// 	config.New(parameters)

// 	// show config content
// 	Print(config)

// 	// start parsing article
// 	articles, err := StartParsing(config)
// 	if err != nil {
// 		test.Error(err)
// 	}
// }