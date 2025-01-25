package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/njayp/theseus/pkg/manager"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) AddImage(config manager.Config) error {
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/add", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send add request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s: %s", resp.Status, body)
	}

	return nil
}

func (c *Client) UpgradeImage(config manager.BuildPayload) error {
	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/upgrade", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send upgrade request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s: %s", resp.Status, body)
	}

	return nil
}

func (c *Client) RemoveImage(imageName string) error {
	data := manager.RemoveRequest{
		ImageName: imageName,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/remove", c.BaseURL), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create remove request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.Post(fmt.Sprintf("%s/remove", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send remove request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed: %s: %s", resp.Status, body)
	}

	return nil
}
