package scholarParser

import (
	"github.com/akmubi/soup"
	"math/rand"
	"time"
	"fmt"
	"log"
)

type Article struct {
	ID int64
	Title string
	URL string
	PDFLink string
	Authors string
	Description string
}

func timeout() {
	time.Sleep( time.Duration(rand.Int() % 8) * time.Second )
}

// returns slice of articles to gived config
func StartParsing(config Config) (articles []Article) {

	var articleCount int

	var lastResponse string

	if err := config.validate(); err != nil {
		log.Fatal(err)
	}

	// 0 - skip first 0 pages
	// 10 - skip first 10 pages
	// ...
	skipSlice := makeIntStringSlice(config.numOfPages)

	// "sample search query?" -> "sample%20search%20query?"
	processedSearchQuery := replaceNonLetters(config.searchQuery)

	fmt.Printf("Searching '%s'...\n", config.searchQuery)

	for _, startValue := range skipSlice {

		fmt.Println("Pages passed:", startValue)

		timeout()

		// building query
		// http://scholar.google.com/scholar/?q=my%20query&start=0&hl=ru
		query := scholarDomain +	"?q=" + processedSearchQuery + 
									"&start=" + startValue +
									"&hl=" + config.language

		fmt.Printf("Query: %s\n\n", query)

		// send request and get a response
		lastResponse, err := soup.Get(query)
		if err != nil {
			log.Fatal(err)
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
					// gs_or_ggsm	(PDF file URL)
					// gs_a			(authors)
					// gs_rs		(description)
					for _, subDiv := range subDivs {
						if subDivClass, hasSubDivClass := subDiv.Attrs()["class"]; hasSubDivClass {

							// title text and URL
							if subDivClass == "gs_ri" {

								textAndURL := subDiv.Find("h3")
								if textAndURLClass, exists := textAndURL.Attrs()["class"]; exists {

									if textAndURLClass == "gs_rt" {

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
									} // titleAndURLClass == "gs_rt"
								} // exists := titleAndURLAttrs()["class"]
							} // title and URL (gs_ri)

							// PDF file URL
							if subDivClass == "gs_or_ggsm" {
								PDFLink := subDiv.Find("a")

								if PDFLink.NodeValue != "" {
									if link, exists := PDFLink.Attrs()["href"]; exists {
										article.PDFLink = link
									}
								}
							}

							// authors
							if config.parseAuthors && subDivClass == "gs_a" {
								article.Authors = subDiv.FullText()
							}

							// description
							if config.parseDescriptions && subDivClass == "gs_rs" {
								article.Description = subDiv.FullText()
							}

						} // hasSubDivClass := subDiv.Attrs()["class"]; hasSubDivClass 
					} // subDiv ~ subDivs

					// END OF ARTICLE SECTION
					article.ID = int64(articleCount)
					articleCount++
					articles = append(articles, article)
				} // divClass == "gs_r gs_or gs_scl"
			} // hasDivClass := div.Attrs()["class"]
		} // div ~ divs
	} // startValue ~ skipSlice

	fmt.Println("Found", articleCount, "pages")

	if articleCount == 0 {
		if lastResponse == "" {
			log.Fatal("No response")
		} else {
			log.Fatal(lastResponse)
		}
	}
	return
}

