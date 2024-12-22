package main

import (
	"fmt"

	"github.com/gammazero/deque"
	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

/**
urls = []

1 - Grab url
2 - find all <a> tags on page
3 - Validate new href="" match conditions
4 - add acceptable urls to urls to process
5 - mark current url as seen/processed

*/

func GetHTML(url string) (string, error) {
	response, err := utils.Get(url)
	if err != nil {
		return "", err
	}
	html := string(response)
	return html, nil
}

func ParseLinks(html string) []string {
	return nil
}

func main() {

	var urlsQ deque.Deque[string]
	urlsQ.PushBack("https://japan-job-en.rakuten.careers/search-jobs")

	for urlsQ.Len() != 0 {
		nextUrl := urlsQ.PopFront()
		html, err := GetHTML(nextUrl)
		if err != nil {
			fmt.Printf("Skipping: %s\n", nextUrl)
		}
		newLinks := ParseLinks(html)
		for _, a := range newLinks {
			urlsQ.PushBack(a)
		}
	}

	// const a_tag_regx = `<a[^>]+href="([^"]+)"`
	// a_rgx := regexp.MustCompile(a_tag_regx)
	// for _, i := range a_rgx.FindAllStringSubmatch(html, -1) {
	// 	fmt.Println(i[1])
	// }

}
