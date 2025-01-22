package manager

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Config struct {
	ContainerConfig *container.Config         `json:"container_config"`
	HostConfig      *container.HostConfig     `json:"host_config"`
	NetworkConfig   *network.NetworkingConfig `json:"network_config"`
}

type RemoveRequest struct {
	ImageName string `json:"image_name"`
}
