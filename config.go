package scholarScraper

import (
	"strconv"
	"errors"
)

type Config struct {
	language	string
	searchQuery	string
	numOfPages	int64
}

const (
	scholarDomain			= "http://scholar.google.com/scholar"
	allowedLanguages		= "ru&en"
	defaultLanguage			= "ru"
	defaultNumOfPages		= 5
)

func (config *Config) New(parameters map[string]string) error {

	// setting up default config parameters
	config.numOfPages	= defaultNumOfPages
	config.language		= defaultLanguage

	if _, has_query := parameters["query"]; !has_query {
		return errors.New("Search query is required!")
	}

	for key, value := range parameters {
		switch key {
		case "lang":
				config.language = value
		case "query":
				config.searchQuery = value
		case "pages":
			numOfPages, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			config.numOfPages = numOfPages
		default:
			return errors.New("Unknown parameter")
		}
	}

	if err := config.validate(); err != nil {
		return err
	}
	return nil
}

func (config Config) validate() error {
	if !isAllowedLanguage(allowedLanguages, config.language) {
		return errors.New("Unknown language")
	}

	// because if page count is greater than 98 a server error occurs
	if config.numOfPages < 0 || config.numOfPages > 98 {
		return errors.New("Incorrect page count")
	}

	if config.searchQuery == "" {
		return errors.New("Seach query is empty")
	}
	return nil
}