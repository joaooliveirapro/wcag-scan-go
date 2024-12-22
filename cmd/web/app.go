package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gammazero/deque"
)

type App struct {
	Workers      int
	MaxDepth     int
	Domain       string
	URLsQ        deque.Deque[string]
	Requests     int
	SeenURLs     sync.Map
	StartingURLs []string
	ExcludeRegex []string
	IncludeRegex []string
	Wg           sync.WaitGroup
	Mut          sync.Mutex
}

func (app *App) Run() {
	// Add starting urls to deque
	for _, url := range app.StartingURLs {
		app.URLsQ.PushBack(url)
	}

	// Start workers
	for workerID := range app.Workers {
		app.Wg.Add(1) // Add Go routine to Waitgroup
		go Trawl(workerID, app)
	}
	// Wait for all Go routines to finish
	app.Wg.Wait()
}

func (app *App) SaveToFile(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()
	for url := range app.SeenURLs.Range {
		if _, err := file.WriteString(fmt.Sprintf("%s\n", url)); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}
}
