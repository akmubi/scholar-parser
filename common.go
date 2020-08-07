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
	fmt.Println("\tLanguage:"			+ config.language + ",")
	fmt.Println("\tSearch Query: '"		+ config.searchQuery + "\",")
	fmt.Println("\tNumber Of Pages:",	config.numOfPages, ",")
	fmt.Println("}")
}

func (article Article) showContent() {
	fmt.Println("Article: {")
	fmt.Println("\tTitle:" 		+ article.Title + ",")
	fmt.Println("\tHTML Link:"	+ article.HTMLLink + ",")
	fmt.Println("\tPDF Link:"	+ article.PDFLink + ",")
	fmt.Println("}")
}
