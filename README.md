# Scholar Parser

Simple Go project for searching papers on Google Scholar, extracting meta and downloading PDF versions of papers.

## Installation

To install the package use following command:

`go get github.com/akmubi/scholar-parser`

(It will be more convenient i you use go modules. Please)

## Function List
```golang
// initialized configuration with given parameters
func (config *Config) New(parameters map[string]string) {
}

// starts parsing Scholar with given configuration
// and returns slice of articles
func StartParsing(config Config) (articles []Article) {
}

// prints content of config, article or link
func Print(object printable) {
}

// downlaods all available PDF versions of articles, 
// saves them into specified location and creates links (article ID ~ PDF filepath) 
func DownloadPDFs(targetFolderPath string, articles []Article) ([]Link, error) {
}

// just saves content of some object to json file
func SaveToJSON(filepath string, object interface{}) (error) {
}

// removes PDF files that aren't actually PDF files
// that may happen if website doesn't want to be parsed 
func RemoveFakePDFs(links []Link) error {
}

// just give an example:
// "--@sample string/hello !!!\u32a7--" -> "sample-string-hello"
func MakeStringPretty(source string) (result string) {
}
```

## Usage example
```golang
package main

import (
	parser "github.com/akmubi/scholar-parser"
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Your query
	query := "Reverse compilation techniques"

	// configuration
	var config parser.Config

	// set searching params
	parameters := map[string]string {
		"query"	: query,
		"lang"	: "en",
		"pages" : "1",
		"parseAuthors" : "yes",
		"parseDesc" : "no",
	}
	
	config.New(parameters)

	// show config content
	parser.Print(config)

	// specify a result folder
	folder := "results/"

	//
	// ARTICLES
	//

	// start parsing article
	articles := parser.StartParsing(config)

	// after parsing we need to save
	// articles information
	parser.SaveToJSON(folder + "articles.json", articles)

	//
	// PDF FILES
	//

	// download available pdf files and create links
	links, err := parser.DownloadPDFs(folder, articles)
	check(err)

	// remove fake PDF files (html) and update links
	err = parser.RemoveFakePDFs(links)
	check(err)


	// and save links to JSON file
	err = parser.SaveToJSON(folder + "links.json", links)
	check(err)
}
```
## Results

### Output

```shell
$ go run main.go
Config: {
        Language: en,
        Number Of Pages: 1,
        Search Query: "Reverse compilation techniques",
        Parse Authors: true,
        Parse Descriptions: false
}
Searching 'Reverse compilation techniques'...
Pages passed: 0
Query: http://scholar.google.com/scholar?q=Reverse%20compilation%20techniques&start=0&hl=en

Pages passed: 10
Query: http://scholar.google.com/scholar?q=Reverse%20compilation%20techniques&start=10&hl=en

Found 20 pages
Saving 'results/articles.json'
[Downloading (1/20)]: Reverse-compilation-techniques.pdf
[Downloading (2/20)]: Low-power-architecture-design-and-compilation-techniques-for-highperformance-processors.pdf
Article 2  doesn't have PDF. Skipping
Article 3  doesn't have PDF. Skipping

...

Article 18  doesn't have PDF. Skipping
[Downloading (20/20)]: Design-for-assembly-techniques-in-reverse-engineering-and-redesign.pdf
Checking article 1's PDF file ...

...

First 10 bytes - '<!doctype '. Removing results/Compilation-techniques-for-recognition-of-parallel-processable-tasks-in-arithmetic-expressions.pdf ...
Checking article 10's PDF file ...
Checking article 11's PDF file ...
First 10 bytes - '%!PS-Adobe'. Removing results/What-assembly-language-programmers-get-up-to-Control-flow-challenges-in-reverse-compilation.pdf ...
Checking article 12's PDF file ...
Checking article 16's PDF file ...
First 10 bytes - '<!DOCTYPE '. Removing results/Reverse-engineering.pdf ...
Checking article 20's PDF file ...
First 10 bytes - '<!DOCTYPE '. Removing results/Design-for-assembly-techniques-in-reverse-engineering-and-redesign.pdf ...
Saving 'results/links.json'
```
### articles.json
Actually, there are 20 articles, but I show only first 3.
```json
[
	{
		"ID": 0,
		"Title": "Reverse compilation techniques",
		"URL": "https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.52.6990\u0026rep=rep1\u0026type=pdf",
		"PDFLink": "https://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.52.6990\u0026rep=rep1\u0026type=pdf",
		"Authors": "C Cifuentes - 1994 - Citeseer",
		"Description": ""
	},
	{
		"ID": 1,
		"Title": "Low power architecture design and compilation techniques for high-performance processors",
		"URL": "https://ieeexplore.ieee.org/abstract/document/282878/",
		"PDFLink": "http://www.scarpaz.com/2100-papers/power%20estimation/su94-low%20power%20architecture%20and%20compilation.pdf",
		"Authors": "CL Su, CY Tsui, AM Despain\ufffd- Proceedings of COMPCON'94, 1994 - ieeexplore.ieee.org",
		"Description": ""
	},
	{
		"ID": 2,
		"Title": "VLIW compilation techniques in a superscalar environment",
		"URL": "https://dl.acm.org/doi/abs/10.1145/178243.178247",
		"PDFLink": "",
		"Authors": "K Ebcioglu, RD Groves, KC Kim…\ufffd- Proceedings of the ACM\ufffd…, 1994 - dl.acm.org",
		"Description": ""
	}
]
```

### links.json
Also only first 3 records.
```json
[
	{
		"ArticleID": 0,
		"Filepath": "results/Reverse-compilation-techniques.pdf"
	},
	{
		"ArticleID": 1,
		"Filepath": "results/Low-power-architecture-design-and-compilation-techniques-for-highperformance-processors.pdf"
	},
	{
		"ArticleID": 4,
		"Filepath": "results/Efficient-compilation-techniques-for-large-scale-feature-models.pdf"
	}
]
```
