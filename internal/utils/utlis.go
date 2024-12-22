package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	LOGS_FOLDER_PATH = "./logs/"
)

func LogAsJSONString(params map[string]any) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
	} else {
		log.Println(string(jsonData))
	}
}

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

func GetHTML(url string) (string, error) {
	// Create a request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		e := fmt.Errorf("[debug] error creating request: %v", err)
		return "", e
	}

	// Create a new HTTP Client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e := fmt.Errorf("[debug] error making request: %v", err)
		return "", e
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Errorf("[debug] error reading response body: %v", err)
		return "", e
	}

	// Response is OK
	if resp.StatusCode != http.StatusOK {
		e := fmt.Errorf("[debug] HTTP code error %d", resp.StatusCode)
		return "", e
	}
	html := string(body)
	return html, nil
}
