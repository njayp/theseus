package manager

import (
	"context"
	"fmt"
	"sync"

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
	ImageName   string
	ImageId     string
	ContainerID string
	Config      Config
}

// NewManager creates a new ImageManager instance
func NewManager() *Manager {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	m := &Manager{
		client: cli,
	}

	err = m.readMap()
	if err != nil {
		m.images = make(map[string]*ImageContainer)
	}

	return m
}

// AddImage adds a new image to manage and ensures its container is running
func (m *Manager) AddImage(ctx context.Context, config Config) error {
	m.Lock()
	defer m.Unlock()

	imageName := config.ContainerConfig.Image

	// Check if image already exists
	if _, exists := m.images[imageName]; exists {
		return fmt.Errorf("image %s is already being managed", imageName)
	}

	// Pull the latest image
	err := m.pullImage(ctx, imageName)
	if err != nil {
		return fmt.Errorf("failed to pull latest image %s: %v", imageName, err)
	}

	digest, _, err := m.client.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return fmt.Errorf("failed to inspect image %s: %v", imageName, err)
	}

	// Create and start container
	containerID, err := m.createAndStartContainer(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create and start container: %v", err)
	}

	// Store the image-container pair
	m.images[imageName] = &ImageContainer{
		ImageName:   imageName,
		ImageId:     digest.ID,
		ContainerID: containerID,
		Config:      config,
	}

	return m.writeMap()
}

// UpgradeImage pulls the latest version of the image and replaces the running container
func (m *Manager) UpgradeImage(ctx context.Context, build BuildPayload) error {
	m.Lock()
	defer m.Unlock()

	imageName := build.Repository.RepoName
	ic, exists := m.images[imageName]
	if !exists {
		return fmt.Errorf("image %s is not being managed", imageName)
	}

	// Pull the latest image
	err := m.pullImage(ctx, imageName)
	if err != nil {
		return fmt.Errorf("failed to pull latest image %s: %v", imageName, err)
	}

	digest, _, err := m.client.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return fmt.Errorf("failed to inspect image %s: %v", imageName, err)
	}

	// If the digest hasn't changed, do nothing
	if digest.ID == ic.ImageId {
		return nil
	}

	// Update the stored image ID
	ic.ImageId = digest.ID

	err = m.stopAndRemoveContainer(ctx, ic.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	// Create and start a new container with the latest image
	newContainerID, err := m.createAndStartContainer(ctx, ic.Config)
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

	if err := m.stopAndRemoveContainer(ctx, ic.ContainerID); err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	delete(m.images, imageName)
	return m.writeMap()
}
