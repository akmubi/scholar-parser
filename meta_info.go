package scholarScraper

import (
	"github.com/akmubi/soup"
	"io/ioutil"
	"os/exec"
	"strings"
	"strconv"
	"errors"
	"fmt"
	// "bytes"
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
	website		string
	journalName	string
	pageCount	int
	number		int
	year		int
}

func (meta pdfMetaInfo) print() {
	fmt.Println("pdfMetaInfo {")
	fmt.Println("\tauthors : [")
	for _, author := range meta.authors {
		fmt.Printf("\t\t'%s',\n", author)
	}
	fmt.Println("\t],")
	fmt.Printf("\tConclusion : '%s',\n", meta.conclusion)
	fmt.Printf("\tAbstract : '%s',\n", meta.abstract)
	fmt.Printf("\tTitle : '%s',\n", meta.name)
	fmt.Printf("\tURL : '%s',\n", meta.url)
	fmt.Printf("\tWebsite : '%s',\n", meta.website)
	fmt.Printf("\tJournal : '%s',\n", meta.journalName)
	fmt.Printf("\tPage Count : %d,\n", meta.pageCount)
	fmt.Printf("\tJournal Number: %d,\n", meta.number)
	fmt.Printf("\tYear: %d,\n", meta.year)
	fmt.Println("}")
}

func extractMeta(articles []Article) (metaInfo []pdfMetaInfo, err error) {
	tempTxtFilesPath := "temp/"
	if err := checkAndCreateFolder(tempTxtFilesPath); err != nil {
		return nil, err
	}

	// var txtFileNames []string
	for _, article := range articles {
		var pdfMeta pdfMetaInfo

		pdfMeta.name = strings.Replace(article.pdfFilePath, "_", " ", -1)
		pdfMeta.url = article.URL
		txtFileName := getFileName(article.pdfFilePath) + ".txt"
		
		// extract text from PDF file 
		cmd := exec.Command("poppler/pdftotext.exe", article.pdfFilePath, tempTxtFilesPath + txtFileName)
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		pdfContent, err := ioutil.ReadFile(tempTxtFilesPath + txtFileName)
		if err != nil {
			return nil, err
		}

		articleInfo, err := article.splitInfo()
		if err != nil {
			return nil, err
		}

		pdfMeta.year = int(articleInfo.year)
		pdfMeta.journalName = articleInfo.journalName
		if articleInfo.isJournalNameLonger {
			pdfMeta.journalName += "..."
		}
		pdfMeta.website = articleInfo.website

		// find author's fullname
		for _, authorInitials := range articleInfo.authors {
			authorFullName, err := findAuthorName(string(pdfContent), authorInitials)
			if err != nil {
				return nil, err
			}
			pdfMeta.authors = append(pdfMeta.authors, authorFullName)
		}

		// find abstract
		abstarctSection, err := FindAbstract(tempTxtFilesPath + txtFileName)
		if err != nil {
			return nil, err
		}
		if pdfMeta.abstract, err = abstarctSection.GetContent(); err != nil {
			return nil, err
		}
		pdfMeta.abstract = SplitSection(pdfMeta.abstract)

		// txtFileNames = append(txtFileNames, tempTxtFilesPath + txtFileName)
		metaInfo = append(metaInfo, pdfMeta)
	}

	err = checkAndRemoveFolder(tempTxtFilesPath)
	return metaInfo, err 
}

type publication struct {
	title		string
	startPage	int64
	endPage		int64
	volume		int64
	issue		int64
	year		int64
}

type scienceDirectMetaInfo struct {
	publication		publication
	authors			[]string
	articleName		string
	abstract		string
}

func (s *scienceDirectMetaInfo) extractMeta(article Article) error {
	if article.URL == "" || !strings.Contains(article.URL, "sciencedirect") {
		return errors.New("Incorrect URL")
	}

	articleInfo, err := article.splitInfo()
	if err != nil {
		return err
	}

	s.publication.year = articleInfo.year
	fmt.Println("Year:", s.publication.year)


	response, err := soup.Get(article.URL)
	if err != nil {
		return err
	}

	document := soup.HTMLParse(response)

	articleNode := document.Find("article")
	if articleNode.NodeValue == "" {
		return errors.New("No article from given URL")
	}

	// first div is publication info
	articleDivs := articleNode.FindAll("div")
	for _, articleDiv := range articleDivs {
		if articleDiv.Attrs()["class"] == "Publication" {
			fmt.Println("[In Publication]")
			divs := articleDiv.FindAll("div")
			for _, div := range divs {

				// publication 
				if strings.Contains(div.Attrs()["class"], "publication-volume") {
					fmt.Println("[In publication-volume]")
					divChildren := div.Children()
					for _, child := range divChildren {
						if strings.Contains(child.Attrs()["class"], "publication-title") {
							fmt.Println("[In publication-title]")
							s.publication.title = child.FullText()
							fmt.Println("Title:", s.publication.title)
						}
						if child.Attrs()["class"] == "text-xs" {
							fmt.Println("[In text-xs]")
							publicationInfo := child.FullText()
							if volume := strings.Index(publicationInfo, "Volume"); volume != -1 {
								volume += 8 // len("Volume ") == 8
								s.publication.volume, _ = strconv.ParseInt(getFirstNumber(publicationInfo[volume:]), 10, 64)
								fmt.Println("Volume:", s.publication.volume)
							}
							if issue := strings.Index(publicationInfo, "Issue"); issue != -1 {
								issue += 6 // len("Issue ") == 6
								s.publication.issue, _ = strconv.ParseInt(getFirstNumber(publicationInfo[issue:]), 10, 64)
								fmt.Println("Issue:", s.publication.issue)
							}
							if pages := strings.Index(publicationInfo, "Pages"); pages != -1 {
								pages += 6 // len("Pages ") == 6
								s.publication.startPage, _ = strconv.ParseInt(getFirstNumber(publicationInfo[pages:]), 10, 64)
								fmt.Println("Star Page:", s.publication.startPage)
								endPageIndex := strings.Index(publicationInfo[pages:], "-")
								s.publication.endPage, _ = strconv.ParseInt(getFirstNumber(publicationInfo[endPageIndex + 1:]), 10, 64)
								fmt.Println("End Page:", s.publication.endPage)
							}
						}
					}
					break
				}
			}
			titleSpan := articleDiv.FindNextSibling().Find("span")

			// Article Title
			if titleSpan.Attrs()["class"] == "title-text" {
				s.articleName = titleSpan.FullText()
				fmt.Println("Article Name:", s.articleName)
			}
		}

		// authors
		if articleDiv.Attrs()["class"] == "author-group" {
			fmt.Println("[In author-group]")
			authorGroups := articleDiv.FindAll("a")
			for i, authorGroup := range authorGroups {
				spans := authorGroup.FindAll("span")
				var author string
				for _, span := range spans {
					if span.Attrs()["class"] == "text given-name" {
						author += span.FullText()
					}
					if span.Attrs()["class"] == "text surname" {
						author += " " + span.FullText()
						break
					}
				}
				fmt.Println("Author", i, "Name:", author)
				s.authors = append(s.authors, author)
			}
		}

		// Abstract
		if articleDiv.Attrs()["class"] == "abstract author" {
			fmt.Println("[In abstract author]")
			abstractDivs := articleDiv.FindAll("div")
			for _, abstractDiv := range abstractDivs {
				if strings.Contains(abstractDiv.Attrs()["id"], "abst") {
					s.abstract += abstractDiv.FullText()
					fmt.Println("Abstract:", s.abstract)
				}
			}
		}
	}
	return nil
}