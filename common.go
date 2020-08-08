package scholarScraper

import (
	"fmt"
)

type printer interface {
	showContent()
}

func Print(object printer) {
	object.showContent()
}

func (config Config) showContent() {
	fmt.Println("Config: {")
	fmt.Println("\tLanguage:",			config.language, ",")
	fmt.Println("\tSearch Query: '",	config.searchQuery, "',")
	fmt.Println("\tNumber Of Pages:",	config.numOfPages)
	fmt.Println("}")
}

func (article Article) showContent() {
	fmt.Println("Article: {")
	fmt.Println("\tTitle:",			article.Title, ",")
	fmt.Println("\tIs HTML:",		article.isHTML, ",")
	fmt.Println("\tPDF Link:",		article.PDFLink, ",")
	fmt.Println("\tURL:",			article.URL, ",")
	fmt.Println("\tInformation:",	article.Info)
	fmt.Println("}")
}
