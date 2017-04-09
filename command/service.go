package command

import "github.com/docker/docker/api/types/swarm"

type ServicesProvider interface {
	Service(id string) swarm.Service
	Order() []string
}

type ServiceContainerMatch interface {
	ServiceId() string

	ContainerIds() []string
	Aliases(id string) []string // Match Container ID (string) using aliases ([]string) provided

}

func serviceContainerMatch_FromString(id string, services []swarm.Service) (ServiceContainerMatch, error) {
	return nil, nil
}
