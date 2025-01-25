package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/njayp/theseus/pkg/manager"
)

func TestClient(t *testing.T) {
	s, err := NewServer()
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	testImage := "jmalloc/echo-server"
	url := "http://localhost" // "http://pi.njayp.net"
	port := 8123
	client := NewClient(fmt.Sprintf("%s:%d", url, port))

	go s.Start(port)
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

func getEnv() []string {
	// Load environment variables from .env file
	envFile := ".env"
	file, err := os.Open(envFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a list of environment variables from the .env file
	envVars := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		envVars = append(envVars, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return envVars
}

func TestProd(t *testing.T) {
	url := "http://pi.njayp.net"
	testImage := "njayp/daedalus"
	client := NewClient(url)

	t.Run("TestAdd", func(t *testing.T) {
		config := manager.Config{
			ContainerConfig: &container.Config{
				Env:   getEnv(),
				Image: testImage,
				ExposedPorts: nat.PortSet{
					nat.Port("6969/tcp"): struct{}{},
				},
			},
			HostConfig: &container.HostConfig{
				PortBindings: nat.PortMap{
					nat.Port("6969/tcp"): []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "6969",
						},
					},
				},
				Mounts: []mount.Mount{
					{
						Type:   mount.TypeBind,
						Source: "/home/njayp/RankedStock.json",
						Target: "/RankedStock.json",
					},
				},
			},
		}

		err := client.AddImage(config)
		if err != nil {
			t.Fatalf("t1: %v", err)
		}
	})

	t.Run("TestUpgrade", func(t *testing.T) {
		err := client.UpgradeImage(manager.BuildPayload{
			Repository: manager.Repository{
				RepoName: testImage,
			},
			PushData: manager.PushData{
				Tag: "latest",
			},
		})
		if err != nil {
			t.Fatalf("t2: %v", err)
		}
	})

	t.Run("TestRemove", func(t *testing.T) {
		err := client.RemoveImage(testImage)
		if err != nil {
			t.Fatalf("t3: %v", err)
		}
	})
}
