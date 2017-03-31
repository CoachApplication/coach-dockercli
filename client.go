package dockercli

import (
	docker_client "github.com/docker/docker/client"
)

func DefaultClient() (*docker_client.Client, error) {
	c, err := docker_client.NewEnvClient()
	return c, err
}
