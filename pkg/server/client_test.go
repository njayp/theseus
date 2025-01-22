package server

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/njayp/theseus/pkg/manager"
)

// var client = NewClient("http://pi.njayp.net")
var client = NewClient("http://localhost:8080")

const testImage = "jmalloc/echo-server"

func TestClient(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	go s.Start(8080)
	time.Sleep(1000 * time.Millisecond)

	t.Run("TestAdd", func(t *testing.T) {
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
	})

	t.Run("TestRemove", func(t *testing.T) {
		err := client.RemoveImage(testImage)
		if err != nil {
			t.Fatalf("t2: %v", err)
		}
	})
}
