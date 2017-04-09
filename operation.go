package dockercli

import (
	docker_client "github.com/docker/docker/client"
)

type ClientOperationBase struct {
	client *docker_client.Client
}

func NewClientOperationBase(c *docker_client.Client) *ClientOperationBase {
	return &ClientOperationBase{
		client: c,
	}
}
func NewClientOperationBaseDefault() *ClientOperationBase {
	c, _ := DefaultClient() // @TODO we don't catch this error, but we should
	return &ClientOperationBase{
		client: c,
	}
}

func (cob *ClientOperationBase) DockerClient() *docker_client.Client {
	return cob.client
}
