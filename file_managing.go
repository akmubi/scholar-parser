package scholarScraper

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"log"
	"fmt"
	"os"
	"io"
)

func DownloadPDFs(targetFolderPath string, articles []Article) ([]Article, error) {
	if err := checkAndCreateFolder(targetFolderPath); err != nil {
		return nil, err
	}

	for i, article := range articles {

		// skip if it's a HTML document
		if !article.isHTML {
			if article.PDFLink == "" {
				fmt.Println("Article", i, " doesn't have PDF. Skipping")
				continue
			}

			prettyName := MakeStringPretty(article.Title)
			if prettyName == "" {
				prettyName += strconv.FormatInt(int64(i), 10)
				log.Println("Non printable title name!")
			}
			prettyName += ".pdf"


			fmt.Println("[Downloading (" + strconv.FormatInt(int64(i + 1), 10) + 
										"/" + strconv.FormatInt(int64(len(articles)), 10) +
										 ")]:", prettyName)

			// Ask for document
			response, err := http.Get(article.PDFLink)
			if err != nil {
				log.Println(err)
				continue
			}

			defer response.Body.Close()

			pdfFile, err := os.Create(targetFolderPath + prettyName)
			if err != nil {
				return nil, err
			}

			defer pdfFile.Close()

			if _, err = io.Copy(pdfFile, response.Body); err != nil {
				log.Println(err)
				continue
			}

			article.pdfFilePath = targetFolderPath + prettyName
		}

	}

	// filter downloaded PDFs
	filteredArticles, err := filterPDFs(articles)
	if err != nil {
		return nil, err
	}

	return filteredArticles, nil
}

func checkAndCreateFolder(filepath string) error {
	folder := getFolder(filepath)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return os.Mkdir(folder, os.ModePerm)
	}
	return nil
}

func SaveToJSON(filepath string, object interface{}) error {
	
	// check if folder exists
	err := checkAndCreateFolder(filepath)
	if err != nil {
		return err
	}

	fmt.Printf("Saving '%s'\n", filepath)
	serialized, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		return err
	}

	// saving json file
	err = ioutil.WriteFile(filepath, serialized, 0644)
	if err != nil {
		return err
	}
	return err
}

func isPdfCorrect(article Article) (bool, error) {
	fmt.Printf("Checking PDF '%s'...\n", article.pdfFilePath)
	
	// Read file
	pdfContent, err := ioutil.ReadFile(article.pdfFilePath)
	if err != nil {
		return false, err
	}
	pdfText := string(pdfContent)
	
	// if first bytes of file are %PDF then it's absolutely PDF file
	// else we need to remove this file
	if !strings.HasPrefix(pdfText, "%PDF") {
		return false, nil
	}
	return true, nil
}

func removePDF(article Article) error {
	if err := os.Remove(article.pdfFilePath); err != nil {
		return err
	}
	return nil
}

func filterPDFs(articles []Article) (filteredArticles []Article, err error) {
	for i := 0; i < len(articles); i++ {
		// check PDF
		isPDF, err := isPdfCorrect(articles[i])
		if err != nil {
			return nil, err
		}

		// If it's a correct PDF file append to result slice
		if isPDF {
			filteredArticles = append(filteredArticles, articles[i])
		// else remove PDF file
		} else {
			if err = removePDF(articles[i]); err != nil {
				return nil, err
			}
		}
	}
	return filteredArticles, nil
}
