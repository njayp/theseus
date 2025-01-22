package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/njayp/theseus/pkg/server"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{BaseURL: baseURL}
}

func (c *Client) AddImage(imageName string) error {
	data := server.AddRequest{
		ImageName: imageName,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/add", c.BaseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to send add request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed: %s", resp.Status)
	}

	return nil
}

func (c *Client) RemoveImage(imageName string) error {
	data := server.RemoveRequest{
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
		return fmt.Errorf("failed to remove image: %s", resp.Status)
	}

	return nil
}
