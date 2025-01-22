package manager

import (
	"context"
	"fmt"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

// ImageManager handles the lifecycle of Docker images and their containers
type ImageManager struct {
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

// NewImageManager creates a new ImageManager instance
func NewImageManager() (*ImageManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	return &ImageManager{
		client: cli,
		images: make(map[string]*ImageContainer),
	}, nil
}

// AddImage adds a new image to manage and ensures its container is running
func (im *ImageManager) AddImage(ctx context.Context, imageName, containerName string) error {
	im.Lock()
	defer im.Unlock()

	// Check if image already exists
	if _, exists := im.images[imageName]; exists {
		return fmt.Errorf("image %s is already being managed", imageName)
	}

	// Pull image if not present
	_, err := im.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %v", imageName, err)
	}

	// Create and start container
	containerID, err := im.createAndStartContainer(ctx, imageName, containerName)
	if err != nil {
		return err
	}

	// Store the image-container pair
	im.images[imageName] = &ImageContainer{
		ImageName:     imageName,
		ContainerID:   containerID,
		ContainerName: containerName,
	}

	return nil
}

// createAndStartContainer creates and starts a new container for the given image
func (im *ImageManager) createAndStartContainer(ctx context.Context, imageName, containerName string) (string, error) {
	resp, err := im.client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container from image %s: %v", imageName, err)
	}

	if err := im.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container %s: %v", resp.ID, err)
	}

	return resp.ID, nil
}

// UpgradeImage pulls the latest version of the image and replaces the running container
func (im *ImageManager) UpgradeImage(ctx context.Context, imageName string) error {
	im.Lock()
	defer im.Unlock()

	ic, exists := im.images[imageName]
	if !exists {
		return fmt.Errorf("image %s is not being managed", imageName)
	}

	// Pull the latest version of the image
	_, err := im.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull latest image %s: %v", imageName, err)
	}

	err = im.removeContainer(ctx, ic.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	// Create and start a new container with the latest image
	newContainerID, err := im.createAndStartContainer(ctx, imageName, ic.ContainerName)
	if err != nil {
		return fmt.Errorf("failed to create new container: %v", err)
	}

	// Update the stored container information
	ic.ContainerID = newContainerID
	return nil
}

// RemoveImage removes an image and stops its container
func (im *ImageManager) RemoveImage(ctx context.Context, imageName string) error {
	im.Lock()
	defer im.Unlock()

	ic, exists := im.images[imageName]
	if !exists {
		return fmt.Errorf("image %s is not being managed", imageName)
	}

	if err := im.removeContainer(ctx, ic.ContainerID); err != nil {
		return fmt.Errorf("failed to remove old container %s: %v", ic.ImageName, err)
	}

	delete(im.images, imageName)
	return nil
}

// removeContainer removes a container by its ID
func (im *ImageManager) removeContainer(ctx context.Context, id string) error {
	err := im.client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return err
	}

	return im.client.ContainerRemove(ctx, id, container.RemoveOptions{})
}
