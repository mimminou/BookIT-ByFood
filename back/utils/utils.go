package utils

import (
	"fmt"
	"github.com/mimminou/BookIT-ByFood/back/models"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

//Utility functions

// Writes anything passed to it to query.txt, This is only used for debugging
func LogQuery(query string) {
	file, _ := os.OpenFile("query.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(query)
	defer file.Close()
}

// Verifiy if fields are empty, return empty fields
// Num pages is not checked because I decided to allow it to be empty

func CheckEmptyFields(book models.Book) []string {
	var emptyFields []string
	if book.Title == "" {
		emptyFields = append(emptyFields, "title")
	}
	if book.Author == "" {
		emptyFields = append(emptyFields, "author")
	}
	if book.Pub_Date == "" {
		emptyFields = append(emptyFields, "pub_date")
	}
	return emptyFields
}

func ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// check if url is valid, and if it's http / https protocle
func IsUrl(link string) bool {
	parsed, err := url.Parse(link)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" { // check if it's http protocole, if not (like ftp or smth else) we don't handle it
		return false
	}
	return true
}

func GetCanonicalUrl(link string) (string, error) {
	regex := `(?i)^https?://[^/]+/[^\?]+`
	re := regexp.MustCompile(regex)

	isMatch := re.MatchString(link)
	if !isMatch {
		return "", fmt.Errorf("URL does not have a canonical format")
	}

	canonicalURL := re.FindString(link)
	//remove trailing slash
	canonicalURL = strings.TrimSuffix(canonicalURL, "/")
	return canonicalURL, nil
}

func GetRedirectionUrl(link string) (string, error) {

	// match byfood.com, case insensitive on 'www' only because technically.
	// a subdomain takeover attack can be performed so we don't match wildcard.
	// a safe approach is to have a known list of active subdomains to match
	FullRegex := `(?i)^https?:\/\/(www\.)?byfood\.com(?:\/.*)?$`
	SubRegex := `(?:https?://)(.*)`
	re := regexp.MustCompile(FullRegex)
	subRegex := regexp.MustCompile(SubRegex)

	isMatch := re.MatchString(link)
	if !isMatch {
		return "", fmt.Errorf("URL is not from ByFood Domain")
	}

	subLink := subRegex.FindStringSubmatch(link)

	if len(subLink) < 2 {
		return "", fmt.Errorf("URL not valid")
	}
	link = subLink[1] // get what's after https?://

	if !strings.HasPrefix(link, "www.") {
		link = fmt.Sprintf("www.%s", link) // add www. if it's not there
	}
	processedURL := "https://" + link // add http(s)://
	processedURL = strings.ToLower(processedURL)
	return processedURL, nil
}
