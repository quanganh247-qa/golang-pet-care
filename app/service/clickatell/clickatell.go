package clickatell

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/quanganh247-qa/go-blog-be/app/util"
)

func SendSMS(ctx context.Context, config *util.Config, phone string, message string) error {
	// URL encode the message content
	encodedMessage := url.QueryEscape(message)
	ollamaURL := fmt.Sprintf("%s?apiKey=%s&to=%s&content=%s", config.ClickatellURL, config.ClickatellAPIKey, phone, encodedMessage)

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(ollamaURL)
	if err != nil {
		return fmt.Errorf("cannot send SMS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot send SMS, status code: %d", resp.StatusCode)
	}

	return nil
}
