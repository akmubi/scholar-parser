package scholarScraper

import (
	"testing"
)

// func TestScienceDirect(test *testing.T) {
// 	article := Article {
// 		Title: "An intelligent value stream-based approach to collaboration of food traceability cyber physical system by fog computing",
// 		isHTML : false,
// 		URL: "https://www.sciencedirect.com/science/article/pii/S0956713516303565",
// 		PDFLink: "",
// 		Info: "RY Chen… - Food Control, 2017 - Elsevier",
// 		pdfFilePath: "",
// 	}
// 	expected := scienceDirectMetaInfo {
// 		publication : publication{ title : "Food Control", startPage : 124, endPage : 136, year : 2017, volume : 71, issue : 0 },
// 		authors : []string{ "Rui-Yang Chen" },
// 		articleName : "An intelligent value stream-based approach to collaboration of food traceability cyber physical system by fog computing",
// 		abstract : "Good advanced food traceability systems help to minimize unsafe or poor quality products in food supply chain through value-based process. From the emerging technologies forthcoming for industry automation, future advanced food traceability system must consider not only cyber physical system (CPS) and fog computing but also value-added business in food supply chain. Accordingly, this study presents a novel intelligent value stream-based food traceability cyber physical system approach integrated with enterprise architectures, EPCglobal and value stream mapping method by fog computing network for traceability collaborative efficiency. Furthermore, the proposed intelligent approach explores distributive and central traceable stream mechanism in assessing the most critical traceable events for tracking and tracing process. Successful case study, software system design and implementation demonstrated the performance of the proposed approach. Furthermore, experiment shows the better results obtained after the simulation execution for intelligent predictive algorithm.",
// 	}

// 	var actual scienceDirectMetaInfo

// 	err := actual.extractMeta(article)
// 	if err != nil {
// 		test.Error("An error occured\n", err)
// 		return
// 	}
// 	if actual.publication.title != expected.publication.title {
// 		test.Error("Mismatch\nExpected:", expected.publication.title, "\nGot:", actual.publication.title)
// 	}
// 	if actual.publication.startPage != expected.publication.startPage {
// 		test.Error("Mismatch\nExpected:", expected.publication.startPage, "\nGot:", actual.publication.startPage)
// 	}
// 	if actual.publication.endPage != expected.publication.endPage {
// 		test.Error("Mismatch\nExpected:", expected.publication.endPage, "\nGot:", actual.publication.endPage)
// 	}
// 	if actual.publication.volume != expected.publication.volume {
// 		test.Error("Mismatch\nExpected:", expected.publication.volume, "\nGot:", actual.publication.volume)
// 	}
// 	if actual.publication.issue != expected.publication.issue {
// 		test.Error("Mismatch\nExpected:", expected.publication.issue, "\nGot:", actual.publication.issue)
// 	}
// 	if actual.publication.year != expected.publication.year {
// 		test.Error("Mismatch\nExpected:", expected.publication.year, "\nGot:", actual.publication.year)
// 	}


// 	if actual.articleName != expected.articleName {
// 		test.Error("Mismatch\nExpected:", expected.articleName, "\nGot:", actual.articleName)
// 	}
// 	if actual.abstract != expected.abstract {
// 		test.Error("Mismatch\nExpected:", expected.abstract, "\nGot:", actual.abstract)
// 	}

// 	if len(actual.authors) != len(expected.authors) {
// 		test.Error("Mismatch\nExpected:", len(expected.authors), "\nGot:", len(actual.authors))
// 	}


// 	for i := range actual.authors {
// 		if expected.authors[i] != actual.authors[i] {
// 			test.Error("Mismatch\nExpected:", expected.authors[i], "\nGot:", actual.authors[i])
// 		}
// 	}
// }



func TestExtractMeta(test *testing.T) {
	article := Article {
		Title: "Cyber–physical system security for the electric power grid",
		isHTML : false,
		URL: "https://ieeexplore.ieee.org/abstract/document/6032699/",
		PDFLink: "https://www.public.asu.edu/~yweng2/Tutorial5/pdf/liu_paper_2.pdf",
		Info: "S Sridhar, A Hahn, M Govindarasu… - Proceedings of the IEEE, 2011 - ieeexplore.ieee.org",
		pdfFilePath: "Cyberphysical-system-security-for-the-electric-power-grid.pdf",
	}

	meta, err := extractMeta([]Article{ article })
	if err != nil {
		test.Error(err)
	}
	meta[0].print()
}