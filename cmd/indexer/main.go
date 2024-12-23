package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

type PageInfo struct {
	URL             string         `json:"url"`
	StatusCode      int            `json:"statusCode"`
	Redirected      bool           `json:"redirected"`
	RedirectHistory []string       `json:"redirectHistory"`
	HtmlTags        map[string]int `json:"HtmlTags"`
	ContentTokens   map[string]int `json:"contentTokens"`
}

type Indexer struct {
	URLsFilePath string
	PageInfo     *PageInfo
}

func (i *Indexer) Stats(response *utils.HTTPResponse) {
	// Regex to capture <tagname>(content)<
	r := regexp.MustCompile(`<(?<tagname>\w*)[^>]*>(?<content>[^<]*)`)
	for _, match := range r.FindAllStringSubmatch(response.HTML, -1) {
		if len(match) < 3 { // Skip matched without content
			continue
		}
		tagName := match[1]
		if len(tagName) == 0 { // Skip empty strings
			continue
		}
		i.HtmlTagsFrequency(tagName) // Tag frequency

		content := match[2]
		i.ContentTokensFrequency(content) // Content tokens frequency
	}
}

func (i *Indexer) HtmlTagsFrequency(tagName string) {
	i.PageInfo.HtmlTags[tagName]++
}

func (i *Indexer) ContentTokensFrequency(content string) {
	for _, w := range strings.Split(content, " ") {
		w = strings.ToLower(w)
		w = string(regexp.MustCompile(`[^\w]`).ReplaceAll([]byte(w), []byte("")))
		if len(w) == 0 {
			continue
		}
		// if there's numbers in w, skip
		if regexp.MustCompile(`\d`).Match([]byte(w)) {
			continue
		}
		i.PageInfo.ContentTokens[w]++
	}
}

func (i *Indexer) CalculateTFIDF() {

}

func (i *Indexer) AppendToJson(filename string) {

	fileData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	allPagesInfo := []PageInfo{}
	if len(fileData) > 0 { // Handle empty file case
		if err := json.Unmarshal(fileData, &allPagesInfo); err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}
	}

	allPagesInfo = append(allPagesInfo, *i.PageInfo)

	updatedData, err := json.MarshalIndent(allPagesInfo, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
	}

	err = os.WriteFile(filename, updatedData, 0644)
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
	}
}

func (i *Indexer) IndexPage(url string) {
	// HTTP request to URL
	response, err := utils.GetHTML(url)
	if err != nil {
		log.Printf("[error] Couldn't get response for %s\n", response.URL)
	}
	// Build PageInfo based on response data
	pageInfo := PageInfo{
		URL:             response.URL,
		StatusCode:      response.StatusCode,
		Redirected:      response.Redirected,
		RedirectHistory: response.RedirectHistory,
		HtmlTags:        map[string]int{},
		ContentTokens:   map[string]int{},
	}
	// 'Attach' PageInfo to indexer
	i.PageInfo = &pageInfo
	i.Stats(&response)
	i.AppendToJson("./cmd/indexer/cache/index.json")
}

func (i *Indexer) Start() {
	// Read URLs file
	file, err := os.OpenFile(i.URLsFilePath, os.O_RDONLY, 0644)
	if err != nil {
		log.Printf("Failed to open file: %v\n", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
	// Iterate each URL in file
	for scanner.Scan() {
		url_ := scanner.Text()
		parsedURL, err := url.Parse(url_)
		URLisValid := err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
		if URLisValid {
			i.IndexPage(url_) // Run Indexer for each valid URL
		}
	}
}

func main() {
	indexer := Indexer{URLsFilePath: "zz_urls.txt"}
	indexer.Start()
}
