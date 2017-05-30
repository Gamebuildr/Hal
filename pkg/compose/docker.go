package compose

import (
	"context"

	"fmt"
	"os"

	"github.com/Gamebuildr/Hal/pkg/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Docker is the implementation of the Docker container engine
type Docker struct {
	Client *client.Client
}

// ContainerNotFound error message when container id cannot be found with image name
const ContainerNotFound string = "container id not found"

// RunContainer will run a docker container
func (engine Docker) RunContainer(message string, image string) error {
	ctx := context.Background()
	env := []string{
		fmt.Sprintf("GCLOUD_PROJECT=%s", os.Getenv("GCLOUD_PROJECT")),
		fmt.Sprintf("GCLOUD_SERVICE_KEY=%s", os.Getenv("GCLOUD_SERVICE_KEY")),
		fmt.Sprintf("PAPERTRAIL_ENDPOINT=%s", os.Getenv(config.LogEndpoint)),
		fmt.Sprintf("REGION=%s", os.Getenv(config.Region)),
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", os.Getenv(config.AWSAccessKeyId)),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", os.Getenv(config.AWSAccessKey)),
		fmt.Sprintf("CODE_REPO_STORAGE=%s", os.Getenv(config.CodeRepoStorage)),
		fmt.Sprintf("QUEUE_URL=%s", os.Getenv(config.QueueURL)),
		fmt.Sprintf("MRROBOT_NOTIFICATIONS=%s", os.Getenv(config.MrrobotNotifications)),
		fmt.Sprintf("GAMEBUILDR_NOTIFICATIONS=%s", os.Getenv(config.GamebuildrNotifications)),
		fmt.Sprintf("GO_ENV=%s", os.Getenv(config.GoEnv)),
		fmt.Sprintf("BUILD_TARGET_PATH=%s", os.Getenv(config.BuildTargetPath)),
		fmt.Sprintf("BUILD_SOURCE_PATH=%s", os.Getenv(config.BuildSourcePath)),
		fmt.Sprintf("ENGINE_LOG_PATH=%s", os.Getenv(config.EngineLogPath)),
		fmt.Sprintf("MESSAGE_STRING=%s", message),
	}

	resp, err := engine.Client.ContainerCreate(ctx, &container.Config{
		Image: image,
		Env:   env,
	}, nil, nil, "")
	if err != nil {
		return err
	}
	if err := engine.Client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

// ActiveContainers returns the list of running docker containers
func (engine Docker) ActiveContainers() (int, error) {
	count := 0

	containers, err := engine.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return count, err
	}
	return len(containers), nil
}

func (engine Docker) getContainerID(image string) (string, error) {
	containers, err := engine.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return "", err
	}
	for _, container := range containers {
		if image == container.Image {
			return container.ID, nil
		}
	}
	return "", fmt.Errorf("%v: %v", ContainerNotFound, image)
}
