package logUtils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	loggerTypes "github.com/Neurofin/requests-logger/store/types"
)

func PostLog(logInput loggerTypes.PostLogInput) {
	logInputJSON, err := json.Marshal(logInput)
	if err != nil {
		log.Printf("Error marshaling log data: %v", err)
		return
	}

	logServiceURL := os.Getenv("LOG_SERVICE_URL")
	if logServiceURL == "" {
		log.Printf("LOG_SERVICE_URL is not set")
		return
	}

	go func() {
		resp, err := http.Post(logServiceURL+"/log", "application/json", bytes.NewBuffer(logInputJSON))
		if err != nil {
			log.Printf("Error posting log data: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Unexpected response from log service: %s", body)
		}
	}()
}
