package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	LOGS_FOLDER_PATH = "./_logs/"
)

func LoggerInit() (*os.File, error) {
	// Set logger properties
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	filepath := fmt.Sprintf("%sdebug_%s.log", LOGS_FOLDER_PATH, time.Now().Format("02-01-2006"))
	logFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)
	return logFile, nil
}

func Get(url string) ([]byte, error) {
	// Create a request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[debug] error creating request: %v\n", err)
		return nil, err
	}

	// Create a new HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[debug] error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[debug] error reading response body: %v\n", err)
		return nil, err
	}

	// Response is OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("[debug] HTTP code error %d\n", resp.StatusCode)
		return nil, err
	}

	return body, nil

}
