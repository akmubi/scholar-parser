package scholarParser

import (
	"strconv"
	"fmt"
)

// common interface
type printable interface {
	showContent()
}

func Print(object printable) {
	object.showContent()
}

func (config Config) showContent() {
	fmt.Println("Config: {")
	fmt.Println("\tLanguage: " + config.language + ",")
	fmt.Println("\tNumber Of Pages: " + strconv.FormatInt(config.numOfPages, 10) + ",")
	fmt.Println("\tSearch Query: \"" + config.searchQuery + "\",")
	fmt.Println("\tParse Authors: " + strconv.FormatBool(config.parseAuthors) + ",")
	fmt.Println("\tParse Descriptions: " + strconv.FormatBool(config.parseDescriptions))
	fmt.Println("}")
}

func (article Article) showContent() {
	fmt.Println("Article: {")
	fmt.Println("\tID: " + strconv.FormatInt(article.ID, 10) + ",")
	fmt.Println("\tTitle: " + article.Title + ",")
	fmt.Println("\tURL: " + article.URL + ",")
	fmt.Println("\tPDF Link: " + article.PDFLink + ",")
	fmt.Println("\tAuthors: " + article.Authors + ",")
	fmt.Println("\tDescription: " + article.Description)
	fmt.Println("}")
}

func (link Link) showContent() {
	fmt.Println("Link : {")
	fmt.Println("\tArticle ID: " + strconv.FormatInt(link.ArticleID, 10) + ",")
	fmt.Println("\tFile Path: " + link.Filepath)
	fmt.Println("}")
}