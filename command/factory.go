package command

import (
	docker_api_types "github.com/docker/docker/api/types/container"

	"github.com/CoachApplication/api"
	"github.com/CoachApplication/command"
)

type Factory struct {
	ContainerIdentifiers []ContainerIdentifier
}

func (f *Factory) AddIdentifier(contId *ContainerIdentifier) error {
	f.ContainerIdentifiers = append(f.ContainerIdentifiers, contId)
	return nil
}

func (f *Factory) NewCommand(id string, ui api.Ui, privileged bool, image string, volumes []string, links []string) command.Command {
	config := docker_api_types.Config{}
	hostConfig := docker_api_types.HostConfig{}

	for _, volume := range volumes {

	}
	for _, link := range links {

	}

	return NewCommand(
		id,
	).Command()
}

type ContainerIdentifier interface {
	ContainerID(reference string) string
}
