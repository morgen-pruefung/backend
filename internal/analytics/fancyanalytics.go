package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const BaseURL = "https://api.fancyanalytics.net"

var fancyAnalyticsClient = NewClient()

type Client struct {
	projectID string
	apiKey    string
}

func NewClient() *Client {
	projectID := mustGetProjectID()
	if projectID == "" {
		fmt.Println("Could not create analytics client, project ID is missing")
		return nil
	}

	apiKey := mustGetAPIKey()
	if apiKey == "" {
		fmt.Println("Could not create analytics client, API key is missing")
		return nil
	}

	return &Client{
		projectID: projectID,
		apiKey:    apiKey,
	}
}

func (c *Client) SendEvent(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/projects/%s/events", BaseURL, c.projectID), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("AuthorizationToken", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func SendEvent(event Event) {
	if fancyAnalyticsClient == nil {
		return
	}

	err := fancyAnalyticsClient.SendEvent(event)
	if err != nil {
		log.Printf("Failed to track event: %v", err)
	}
}

func mustGetProjectID() string {
	projectID := os.Getenv("FANCYANALYTICS_PROJECT_ID")
	if projectID == "" {
		return ""
	}

	return projectID
}

func mustGetAPIKey() string {
	apiKey := os.Getenv("FANCYANALYTICS_API_KEY")
	if apiKey == "" {
		return ""
	}

	return apiKey
}
