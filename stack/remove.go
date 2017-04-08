package stack

import (
	"context"
	"fmt"

	docker_api_types "github.com/docker/docker/api/types"
	docker_api_types_swarm "github.com/docker/docker/api/types/swarm"
	docker_cli_command "github.com/docker/docker/cli/command"
)

func removeServices(
	ctx context.Context,
	dockerCli docker_cli_command.Cli,
	services []docker_api_types_swarm.Service,
) bool {
	var err error
	for _, service := range services {
		fmt.Fprintf(dockerCli.Err(), "Removing service %s\n", service.Spec.Name)
		if err = dockerCli.Client().ServiceRemove(ctx, service.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove service %s: %s", service.ID, err)
		}
	}
	return err != nil
}

func removeNetworks(
	ctx context.Context,
	dockerCli docker_cli_command.Cli,
	networks []docker_api_types.NetworkResource,
) bool {
	var err error
	for _, network := range networks {
		fmt.Fprintf(dockerCli.Err(), "Removing network %s\n", network.Name)
		if err = dockerCli.Client().NetworkRemove(ctx, network.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove network %s: %s", network.ID, err)
		}
	}
	return err != nil
}

func removeSecrets(
	ctx context.Context,
	dockerCli docker_cli_command.Cli,
	secrets []docker_api_types_swarm.Secret,
) bool {
	var err error
	for _, secret := range secrets {
		fmt.Fprintf(dockerCli.Err(), "Removing secret %s\n", secret.Spec.Name)
		if err = dockerCli.Client().SecretRemove(ctx, secret.ID); err != nil {
			fmt.Fprintf(dockerCli.Err(), "Failed to remove secret %s: %s", secret.ID, err)
		}
	}
	return err != nil
}
