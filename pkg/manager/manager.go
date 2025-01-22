package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// Manager handles the lifecycle of Docker images and their containers
type Manager struct {
	sync.RWMutex
	client *client.Client
	images map[string]*ImageContainer
}

// ImageContainer represents a Docker image and its running container
type ImageContainer struct {
	ImageName     string
	ContainerID   string
	ContainerName string
}

// NewManager creates a new ImageManager instance
func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &Manager{
		client: cli,
		images: make(map[string]*ImageContainer),
	}, nil
}

// AddImage adds a new image to manage and ensures its container is running
func (m *Manager) AddImage(ctx context.Context, imageName, containerName string) error {
	m.Lock()
	defer m.Unlock()

	// Check if image already exists
	if _, exists := m.images[imageName]; exists {
		return fmt.Errorf("image %s is already being managed", imageName)
	}

	// Pull image if not present
	_, err := m.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %v", imageName, err)
	}

	// Create and start container
	containerID, err := m.createAndStartContainer(ctx, imageName, containerName)
	if err != nil {
		return err
	}

	// Store the image-container pair
	m.images[imageName] = &ImageContainer{
		ImageName:     imageName,
		ContainerID:   containerID,
		ContainerName: containerName,
	}

	return nil
}

// createAndStartContainer creates and starts a new container for the given image
func (m *Manager) createAndStartContainer(ctx context.Context, imageName, containerName string) (string, error) {
	resp, err := m.client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container from image %s: %v", imageName, err)
	}

	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container %s: %v", resp.ID, err)
	}

	return resp.ID, nil
}

// UpgradeImage pulls the latest version of the image and replaces the running container
func (m *Manager) UpgradeImage(ctx context.Context, imageName string) error {
	m.Lock()
	defer m.Unlock()

	ic, exists := m.images[imageName]
	if !exists {
		return fmt.Errorf("image %s is not being managed", imageName)
	}

	// Pull the latest version of the image
	_, err := m.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull latest image %s: %v", imageName, err)
	}

	err = m.removeContainer(ctx, ic.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	// Create and start a new container with the latest image
	newContainerID, err := m.createAndStartContainer(ctx, imageName, ic.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to create new container: %v", err)
	}

	// Update the stored container information
	ic.ContainerID = newContainerID
	return nil
}

// RemoveImage removes an image and stops its container
func (m *Manager) RemoveImage(ctx context.Context, imageName string) error {
	m.Lock()
	defer m.Unlock()

	ic, exists := m.images[imageName]
	if !exists {
		return fmt.Errorf("image %s is not being managed", imageName)
	}

	if err := m.removeContainer(ctx, ic.ContainerID); err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	delete(m.images, imageName)
	return nil
}

// removeContainer removes a container by its ID
func (m *Manager) removeContainer(ctx context.Context, id string) error {
	err := m.client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return err
	}

	return m.client.ContainerRemove(ctx, id, container.RemoveOptions{})
}
