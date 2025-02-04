package manager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/njayp/theseus/pkg/util"
)

// createAndStartContainer creates and starts a new container for the given image
func (m *Manager) createAndStartContainer(ctx context.Context, config Config) (string, error) {
	resp, err := m.client.ContainerCreate(ctx, config.ContainerConfig, config.HostConfig, config.NetworkConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container from image %s: %v", config.ContainerConfig.Image, err)
	}

	if err := m.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container %s: %v", resp.ID, err)
	}

	return resp.ID, nil
}

// stopAndRemoveContainer removes a container by its ID
func (m *Manager) stopAndRemoveContainer(ctx context.Context, id string) error {
	err := m.client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		return err
	}

	return m.client.ContainerRemove(ctx, id, container.RemoveOptions{})
}

func (m *Manager) pullImage(ctx context.Context, imageName string) error {
	// Pull the latest version of the image
	reader, err := m.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull latest image %s: %v", imageName, err)
	}
	defer reader.Close()

	// Print the pull output
	_, err = io.Copy(os.Stdout, reader)
	return err
}

const filename = "/mnt/map.json"

func (m *Manager) writeMap() error {
	return util.WriteJson(filename, m.images)
}

func (m *Manager) readMap() error {
	images, err := util.ReadJson[map[string]*ImageContainer](filename)
	if err != nil {
		return err
	}

	m.images = images
	return nil
}
