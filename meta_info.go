package scholarScraper

import (
	"strings"
	"errors"
	"github.com/akmubi/soup"
)

type metaExtractor interface {
	extractMeta(article Article) error
}

type htmlMetaInfo struct {
	authors		[]string
	webSite		string
	name		string
	url			string
}

type pdfMetaInfo struct {
	authors		[]string
	conclusion	string
	abstract	string
	name		string
	url			string
	pageCount	int
	number		int
}

func (pdfInfo *pdfMetaInfo) extractMeta(article Article) error {

}

// type publication struct {
// 	volume		string
// 	issue		string
// 	name		string
// 	startPage	int
// 	endPage		int
// 	year		int
// }

// type scienceDirectMetaInfo struct {
// 	publication		publication
// 	authors			[]string
// 	acticleName		string
// 	abstract		string
// }

// func (s *scienceDirectMetaInfo) extractMeta(article Article) error {
// 	if article.URL == "" || !strings.Contains(article.URL, "sciencedirect") {
// 		return errors.New("Incorrect URL")
// 	}

// 	response, err := soup.Get(article.URL)
// 	if err != nil {
// 		return err
// 	}

// 	document := soup.HTMLParse(response)

// 	article := document.Find("article")
// 	if article.NodeValue == "" {
// 		return errors.New("No article from given URL")
// 	}

// }