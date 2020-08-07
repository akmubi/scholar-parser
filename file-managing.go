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

// Article Info ~ PDF File
type Link struct {
	ArticleID int64
	Filepath string
}

func findFilenameAndNull(filepath string, links []Link) {
	for _, link := range links {
		if link.Filepath == filepath {
			link.Filepath = ""
			return
		}
	}
} 

// Serialization

func DownloadPDFs(targetFolderPath string, articles []Article) ([]Link, error) {

	var links []Link

	err := checkAndCreateFolder(targetFolderPath)
	if err != nil {
		return nil, err
	}

	for i, article := range articles {

		var link Link

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

		link.ArticleID = article.ID

		// Ask for document
		response, err := http.Get(article.PDFLink)
		if err != nil {
			// do not return error
			// because we need to download
			// as much as possible
			log.Println(err)
			continue
		}
		defer response.Body.Close()

		pdfFile, err := os.Create(targetFolderPath + prettyName)
		if err != nil {
			return nil, err
		}

		defer pdfFile.Close()

		_, err = io.Copy(pdfFile, response.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		link.Filepath = targetFolderPath + prettyName
		links = append(links, link)
	}
	return links, nil
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

func (link *Link) checkAndRemovePDF() error {
	fmt.Printf("Checking article %d's PDF file ...\n", link.ArticleID + 1)
	
	// Read file
	PDFContent, err := ioutil.ReadFile(link.Filepath)
	if err != nil {
		return err
	}
	PDFContentString := string(PDFContent)
	
	// if first bytes of file are %PDF then it's absolutely PDF file
	// else we need to remove this file
	if !strings.HasPrefix(PDFContentString, "%PDF") {
		fmt.Printf("First 10 bytes - '%s'. Removing %s ...\n", PDFContentString[:10], link.Filepath)

		err = os.Remove(link.Filepath)
		
		if err != nil {
			return err
		}

		// Remove filepath from link
		link.Filepath = ""
	}
	return nil
}

func RemoveFakePDFs(links []Link) error {
	for _, link := range links {
		err := link.checkAndRemovePDF()
		if err != nil {
			return err
		}
	}
	return nil
}