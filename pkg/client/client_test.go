package client

import (
	"testing"
)

func TestClient(t *testing.T) {
	// Create a new client
	client := NewClient("http://localhost:8080")

	// Test AddImage
	err := client.AddImage("jmalloc/echo-server")
	if err != nil {
		t.Fatalf("t1: %v", err)
	}

	// Test RemoveImage
	err = client.RemoveImage("jmalloc/echo-server")
	if err != nil {
		t.Fatalf("t2: %v", err)
	}
}
