package manager

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestManager(t *testing.T) {
	ctx := context.Background()

	// Create a new manager
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Define test image and container names
	imageName := "jmalloc/echo-server"
	containerName := "test-container"

	// Test AddImage
	err = mgr.AddImage(ctx, imageName, containerName)
	if err != nil {
		t.Fatalf("Failed to add image: %v", err)
	}

	// Verify container is running
	containers, err := mgr.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	found := false
	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Container %s not found", containerName)
	}

	// Test UpgradeImage
	err = mgr.UpgradeImage(ctx, imageName)
	if err != nil {
		t.Fatalf("Failed to upgrade image: %v", err)
	}

	// Test RemoveImage
	err = mgr.RemoveImage(ctx, imageName)
	if err != nil {
		t.Fatalf("Failed to remove image: %v", err)
	}

	// Verify container is removed
	containers, err = mgr.client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}

	for _, container := range containers {
		if container.Names[0] == "/"+containerName {
			t.Fatalf("Container %s was not removed", containerName)
		}
	}
}
