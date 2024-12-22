package main

import (
	"fmt"

	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

func main() {
	// Logger init
	log, err := utils.LoggerInit()
	if err != nil {
		fmt.Printf("[fatal] Logger didn't start: %+v", err)
	}
	defer log.Close()

	// App config
	// For each worker, must add a starting url to prevent early closure of workers
	app := App{
		Workers:  2,
		MaxDepth: 40,
		Domain:   "careers.arm.com",
		StartingURLs: []string{
			"https://careers.arm.com/",
			"https://careers.arm.com/search-jobs",
		},
		ExcludeRegex: []string{},
		IncludeRegex: []string{`/job/`},
	}

	// Run app
	app.Run()

	// Save seen URLs to file
	app.SaveToFile("seen.txt")

	//
	fmt.Printf("%d\n", app.Requests)
}
