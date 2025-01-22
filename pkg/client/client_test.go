package client

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/njayp/theseus/pkg/manager"
)

// var client = NewClient("http://pi.njayp.net")
var client = NewClient("http://localhost:8080")

const testImage = "jmalloc/echo-server"

func TestAdd(t *testing.T) {
	config := manager.Config{
		ContainerConfig: &container.Config{
			Image: testImage,
			ExposedPorts: nat.PortSet{
				nat.Port("8080/tcp"): struct{}{},
			},
		},
		HostConfig: &container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port("8080/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "80",
					},
				},
			},
		},
	}

	err := client.AddImage(config)
	if err != nil {
		t.Fatalf("t1: %v", err)
	}
}

func TestRemove(t *testing.T) {
	err := client.RemoveImage(testImage)
	if err != nil {
		t.Fatalf("t2: %v", err)
	}
}
