package main

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

func FindAllPageLinks(html string) ([]string, error) {
	newLinks := []string{}
	a_tag_regex := regexp.MustCompile(`<a[^>]*href=["']([^"']+)["']`)
	matches := a_tag_regex.FindAllStringSubmatch(html, -1)
	if len(matches) > 1 {
		// Found a tags with href=""
		for _, a := range matches {
			// a[0] -> (<a ... ) all a tag up to ">"
			// a[1] -> href="(this part)"
			newLinks = append(newLinks, a[1])
		}
	}
	if len(newLinks) == 0 {
		return newLinks, fmt.Errorf("no <a> tags found")
	}
	return newLinks, nil
}

func ValidateLink(a string, domain string) bool {
	// Ignore "#..." anchor links
	if strings.HasPrefix(a, "#") {
		return false
	}
	// Ignore mailto: links
	if strings.Contains(a, "mailto:") {
		return false
	}
	// Ignore "/" links to main page
	if a == "/" {
		return false
	}
	// Ignore links that don't contain the domain
	if strings.HasPrefix(a, "http") && !strings.Contains(a, domain) {
		return false
	}
	// Ignore social media links that contain the domain link
	if strings.Contains(a, domain) {
		return strings.HasPrefix(a, fmt.Sprintf("https://%s", domain))
	}
	return true
}

func NormaliseLink(a string, domain string) (string, error) {
	var rawUrl = a
	// Ignore // links
	if strings.HasPrefix(rawUrl, "//") {
		return a, nil
	}
	// Normalise relative URLs
	if strings.HasPrefix(rawUrl, "/") {
		rawUrl = fmt.Sprintf("https://%s%s", domain, rawUrl)
	}
	// Only process URLs starting with "http"
	if strings.HasPrefix(rawUrl, "http") {
		// Workaround to ensure invalid URLs with HTML entities are parsed correctly
		// Example: /path&amps;p=something
		// This get convert to /path?amps;p=something whis is there removed from the URL
		rawUrl = strings.Replace(rawUrl, "&", "?", 1)
		// Parse URL
		parsedUrl, err := url.Parse(rawUrl)
		if err != nil {
			return "", err
		}
		// Unescape the RawQuery part of the URL
		parsedUrl.RawQuery = html.UnescapeString(parsedUrl.RawQuery)
		// Remove query params
		parsedUrl.RawQuery = ""
		return parsedUrl.String(), nil
	} else {
		return "", fmt.Errorf("invalid link: %s", rawUrl)
	}
}

func Trawl(workerID int, app *App) {
	// Defer signal that go routine is done
	defer app.Wg.Done()

	fmt.Printf("[Worker %d] is starting\n", workerID)

	/* Run while there's URLs to process or stop condition has been reached (safe stop) */
	for {
		/* Check URLsQ size and stop condition. Grab lock */
		app.Mut.Lock()
		noMoreUrls := app.URLsQ.Len() == 0
		stopConditionReached := app.Requests >= app.MaxDepth
		if noMoreUrls || stopConditionReached {
			app.Mut.Unlock() // Release lock and break out of loop
			break
		}
		nextUrl := app.URLsQ.PopFront() // PopFront() needs to be thread safe to avoid panics
		app.Mut.Unlock()                // Release lock after popping the URL

		fmt.Printf("[Worker %d] parsing %s\n", workerID, nextUrl)

		/* Get HTML string for URL */
		html, err := utils.GetHTML(nextUrl)
		if err != nil {
			utils.LogAsJSONString(map[string]any{
				"workerid": workerID,
				"url":      nextUrl,
				"func":     "trawl",
				"error":    fmt.Sprintf("%+v", err.Error()),
			})
		}

		/* Increment Requests count */
		app.Requests++

		/* Find all <a> on page */
		newLinks, err := FindAllPageLinks(html)
		if err != nil {
			utils.LogAsJSONString(map[string]any{
				"workerid": workerID,
				"url":      nextUrl,
				"func":     "trawl",
				"error":    fmt.Sprintf("%+v", err.Error()),
			})
		}

		/* Validate and normalise links found */
		app.Mut.Lock()
		var added = 0
		for _, a := range newLinks {
			normalA, err := NormaliseLink(a, app.Domain)
			if err != nil {
				utils.LogAsJSONString(map[string]any{
					"workerid": workerID,
					"url":      nextUrl,
					"func":     "trawl",
					"anchor":   a,
					"error":    fmt.Sprintf("%+v", err.Error()),
				})
				// Log and skip to next URL
				continue
			}
			rejectUrl := false
			ok := ValidateLink(normalA, app.Domain)
			if ok {
				// Take into consideration URL exclusion regexes
				for _, excludeUrlRegex := range app.ExcludeRegex {
					if regexp.MustCompile(excludeUrlRegex).MatchString(normalA) {
						rejectUrl = true
						break // as soon as one regex is match, the url is no longer valid
					}
				}
				// Take into consideration URL include regexes
				for _, includeUrlRegex := range app.IncludeRegex {
					if !regexp.MustCompile(includeUrlRegex).MatchString(normalA) {
						rejectUrl = true
						break
					}
				}
				if rejectUrl {
					// Skip this url as it's been excluded due to app.ExcludeRegex
					continue
				}
				// sync.Map LOAD or STORE method returns item and boolean if it was already on map
				// Keep track of what urls have been added already to the Q
				_, hasBeenSeen := app.SeenURLs.LoadOrStore(normalA, true)
				if !hasBeenSeen {
					// Thread safe pushing next URL to Deque
					// as Mut is locked outside for loop on line 149
					app.URLsQ.PushBack(normalA)
					added++
				}
			}
		}
		fmt.Printf("[Worker %d] added %d new links to URLsQ.\n", workerID, added)
		app.Mut.Unlock() // Cannot use defer for Unlock() otherwise go routines hang
	}
	fmt.Printf("[Worker %d] is done.\n", workerID)
}
