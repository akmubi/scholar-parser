package scholarScraper

import (
	"github.com/akmubi/soup"
	"math/rand"
	"time"
	"fmt"
	"strconv"
	"strings"
	"errors"
	// "log"
)

type Article struct {
	Title		string
	// is article actually a HTML document
	isHTML		bool
	// link to download the PDF file (if it's available)
	PDFLink		string
	// link to a web-site
	URL			string
	// authors, publisher year, web-site
	Info		string
	pdfFilePath	string
}

func timeout() {
	duration := 10 + rand.Int() % 10
	fmt.Println("Timeout -", duration)
	for i := 0; i < duration; i++ {
		fmt.Printf("\r[%d]", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println()
}

// returns slice of articles to gived config
func StartParsing(config Config) (articles []Article, err error) {

	var articleCount int

	var lastResponse string

	// 0 - skip first 0 pages
	// 10 - skip first 10 pages
	// ...
	if config.numOfPages == 0 {
		fmt.Printf("You put 0 to page count. OK")
		return articles, nil
	}
	skipSlice := makeIntStringSlice(config.numOfPages)

	// "sample search query?" -> "sample%20search%20query?"
	processedSearchQuery := replaceNonLetters(config.searchQuery)

	fmt.Printf("Searching '%s'...\n", config.searchQuery)

	for _, startValue := range skipSlice {

		// building query
		// http://scholar.google.com/scholar/?q=my%20query&start=0&hl=ru
		query := scholarDomain +	"?q=" + processedSearchQuery + 
									"&start=" + startValue +
									"&hl=" + config.language

		fmt.Printf("Query: %s\n\n", query)

		// send request and get a response
		lastResponse, err := soup.Get(query)
		if err != nil {
			return nil, err
		}

		// parsing recieved HTML document
		// and create DOM tree
		document := soup.HTMLParse(lastResponse)

		// find all <div> tags
		divs := document.FindAll("div")

		for _, div := range divs {

			// step through all div class names
			if divClass, hasDivClass := div.Attrs()["class"]; hasDivClass {

				// ARTICLE SECTION (title text, URL, authors, desctiption)
				// <div class="gs_r gs_or gs_scl" ...>
				if divClass == "gs_r gs_or gs_scl" {
					
					// START OF ARTICLE SECTION
					var article Article

					subDivs := div.FindAll("div")

					// looking for:
					// gs_ri		(title text and URL)
					// gs_or_ggsm	(PDF/HTML link)
					// gs_a			(authors)
					for _, subDiv := range subDivs {
						if subDivClass, hasSubDivClass := subDiv.Attrs()["class"]; hasSubDivClass {

							// title text and URL
							if subDivClass == "gs_ri" {

								textAndURL := subDiv.Find("h3")
								if textAndURLClass, exists := textAndURL.Attrs()["class"]; exists {

									if textAndURLClass == "gs_rt" {

										bookPrefix :=  textAndURL.Find("span")
										var bookLabel string
										
										if bookPrefix.NodeValue != "" {
											bookLabel = "[BOOK]"
											if config.language == "ru" {
												bookLabel = "[КНИГА]"
											}

											// we don't need BOOKS, we need articles
											if bookPrefix.FullText() == bookLabel {
												continue
											}
										}

										// find <a> tag
										title := textAndURL.Find("a")

										// check is performing because
										// there is may be a quote insted of
										// title and URL
										if title.NodeValue != "" {

											// title text
											article.Title = unicodeConvert(title.FullText())

											// URL
											if URL, hasURL := title.Attrs()["href"]; hasURL {
												article.URL = URL
											}
										} else {

											// check if it's a quote
											spans := textAndURL.FindAll("span")

											if spans != nil {
												for _, span := range spans {

													// if span has id then we need to get text
													if _, hasId := span.Attrs()["id"]; hasId {
														article.Title = unicodeConvert(span.FullText())
													}
												}
											}
										}
									}
								}
							}

							// PDF/HTML link
							if subDivClass == "gs_or_ggsm" {
								Link := subDiv.Find("a")
								if Link.NodeValue != "" {
									if docType := Link.Find("span"); docType.NodeValue != "" && docType.Attrs()["class"] == "gs_ctg2" {
										if docType.FullText() == "[HTML]" {
											article.isHTML = true
										} else if docType.FullText() == "[PDF]" {
											if pdfLink, exists := Link.Attrs()["href"]; exists {
												article.PDFLink = pdfLink
											}
										} else {
											return nil, errors.New("Unknown link type")
										}
									}
								}
							}

							// authors, publisher, year, web-site name
							if subDivClass == "gs_a" {
								article.Info = subDiv.FullText()
							}
						}
					}

					// END OF ARTICLE SECTION
					articleCount++
					articles = append(articles, article)
				}
			}
		}
		pagesPassed, err := strconv.ParseInt(startValue, 10, 64)
		if err != nil {
			return nil, err
		}

		fmt.Println("Pages passed:", pagesPassed + int64(10))
		timeout()
	} // startValue ~ skipSlice

	fmt.Println("Found", articleCount, "pages")

	if articleCount == 0 {
		if lastResponse == "" {
			return nil, errors.New("No response")
		} else {
			return nil, errors.New(lastResponse)
		}
	}
	return articles, nil
}

type info struct {
	authors				[]string
	journalName			string
	website				string
	isJournalNameLonger bool
	year				int64
}

func emptyInfo() info {
	return info {
		authors : nil,
		journalName : "",
		website : "",
		isJournalNameLonger : false,
		year : -1,
	}
}

// Example: "AA Author1, B Author2... - 2018 - dl.acm.org" ==>
// 			"AA Author1", "B Author2", 2018, dl.acm.org
func (article Article) splitInfo() (articleInfo info, err error) {
	parts := strings.Split(article.Info, " - ")
	partsLength := len(parts)

	var website, journal, year, authors string
	if partsLength == 4 {
		website = parts[3]
		year	= parts[2]
		journal = parts[1]
		authors = parts[0]
	} else if partsLength == 3 {
		website = parts[2]
		authors = parts[0]

		// try to find year
		year = parts[1]
		year = year[len(year) - 4:]
		_, err := strconv.ParseInt(year, 10, 64)
		if err != nil {
			year = ""
		}

		journal = parts[1]
		journal = journal[:len(journal) - 4]

	} else {
		return emptyInfo(), errors.New("Unknown article info form")
	}

	articleInfo.website = strings.Trim(website, " ")
	articleInfo.year, err = strconv.ParseInt(year, 10, 64)
	if err != nil {
		return emptyInfo(), err
	}

	articleInfo.isJournalNameLonger = strings.ContainsRune(journal, '…')
	articleInfo.journalName = strings.Trim(journal, " ,…")

	// authors
	authors = strings.Trim(authors, " …")
	articleInfo.authors = strings.Split(authors, ", ")
	return articleInfo, nil
}