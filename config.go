package scholarParser

import (
	// "strings"
	"errors"
	"log"
	"strconv"
)

type Config struct {
	language          string
	numOfPages        int64
	searchQuery       string
	parseAuthors      bool
	parseDescriptions bool
}

const (
	scholarDomain            = "http://scholar.google.com/scholar"
	defaultLanguage          = "ru"
	defaultNumOfPages        = 5
	defaultParseAuthors      = false
	defaultParseDescriptions = false
	allowedLanguages         = "ru&en"
)

func (config *Config) New(parameters map[string]string) {

	// setting up default config parameters
	config.language = defaultLanguage
	config.numOfPages = defaultNumOfPages
	config.parseAuthors = defaultParseAuthors
	config.parseDescriptions = defaultParseDescriptions

	if _, has_query := parameters["query"]; !has_query {
		log.Fatal("Search query is required!")
	}

	for key, value := range parameters {
		switch key {
		case "lang":
			if isAllowedLanguage(allowedLanguages, value) {
				config.language = value
			} else {
				log.Fatalf("Unknown '%s' parameter value (%s). Must be 'ru/en'", key, value)
			}
		case "query":
			if value != "" {
				config.searchQuery = value
			} else {
				log.Fatalln("Search query is not specified")
			}
		case "pages":
			numOfPages, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				log.Fatalln(err)
			}
			config.numOfPages = numOfPages
		case "parseAuthors":
			if value != "" {
				if value == "yes" || value == "true" {
					config.parseAuthors = true
				} else if value == "no" || value == "false" {
					config.parseAuthors = false
				} else {
					log.Fatalf("Unknown '%s' parameter value (%s). Must be - 'yes/true/no/false'", key, value)
				}
			} else {
				log.Fatalf("Parameter '%s' is not specified\n", key)
			}
		case "parseDesc":
			if value != "" {
				if value == "yes" || value == "true" {
					config.parseDescriptions = true
				} else if value == "no" || value == "false" {
					config.parseDescriptions = false
				} else {
					log.Fatalf("Unknown '%s' parameter value (%s). Must be - 'yes/true/no/false'", key, value)
				}
			} else {
				log.Fatalf("Parameter '%s' is not specified\n", key)
			}
		default:
			log.Fatalf("Unknown parameter - '%s'\n", key)
		}
	}
}

func (config Config) validate() error {
	if config.searchQuery == "" {
		return errors.New("Search query is required")
	}

	if !isAllowedLanguage(allowedLanguages, config.language) {
		return errors.New("Unknown language value")
	}

	if config.numOfPages <= 0 {
		return errors.New("Invalid number of pages")
	}

	return nil
}
