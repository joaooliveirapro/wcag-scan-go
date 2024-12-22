package main

import (
	"fmt"
	"os"

	"github.com/joaooliveirapro/wcag-scan-go/internal/utils"
)

func main() {
	// Logger init
	log, err := utils.LoggerInit()
	if err != nil {
		fmt.Printf("[fatal] Logger didn't start: %+v", err)
		os.Exit(1)
	}
	defer log.Close()

	// App config
	// For each worker, must add a starting url to prevent early closure of workers
	app := App{
		Workers:  1,
		MaxDepth: 20,
		Domain:   "careers.adeccogroup.com",
		StartingURLs: []string{
			"https://careers.adeccogroup.com/",
		},
		ExcludeRegex: []string{},
		IncludeRegex: []string{},
	}

	// Run app
	app.Run()

	// Save processed URLs to file
	app.SaveProcessedToFile("urls_processed.txt")

}
