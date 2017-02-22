package compose

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Docker is the implementation of the Docker container engine
type Docker struct {
	client *client.Client
}

// RunContainer will run a docker container
func (engine Docker) RunContainer(image string) error {
	containerID, err := engine.getContainerID(image)
	if err != nil {
		return err
	}
	ctx := context.Background()
	if err := engine.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

// ActiveContainers returns the list of running docker containers
func (engine Docker) ActiveContainers() (int, error) {
	count := 0

	containers, err := engine.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return count, err
	}
	return len(containers), nil
}

func (engine Docker) getContainerID(image string) (string, error) {
	containers, err := engine.client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		if image == container.Image {
			return container.ID, nil
		}
	}
	return "", errors.New("container id not found")
}
