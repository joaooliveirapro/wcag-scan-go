package main

import (
	"fmt"
	"time"

	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

type PageInfo struct {
	HTTPResponse  utils.HTTPResponse `json:"HTTPResponse"`
	HTMLTags      map[string]int     `json:"HtmlTags"`
	ContentTokens map[string]int     `json:"contentTokens"`
	Timestamp     string             `json:"timestamp"`
}

func NewPageInfo(url string) (*PageInfo, error) {
	// HTTP request to URL
	response, err := utils.GetHTML(url)
	if err != nil {
		return nil, fmt.Errorf("[error] Couldn't get response for %s", response.URL)
	}
	// Build PageInfo based on response data
	pageInfo := PageInfo{
		Timestamp:     time.Now().Format("02-01-2006 15:04:05"),
		HTTPResponse:  response,
		HTMLTags:      map[string]int{},
		ContentTokens: map[string]int{},
	}
	return &pageInfo, nil
}
