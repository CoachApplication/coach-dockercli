package stack

import (
	"context"
	"fmt"
	docker_api_types "github.com/docker/docker/api/types"
	docker_api_types_swarm "github.com/docker/docker/api/types/swarm"
	docker_cli_command "github.com/docker/docker/cli/command"
	docker_cli_compose_convert "github.com/docker/docker/cli/compose/convert"
	docker_client "github.com/docker/docker/client"
)

var defaultNetworkDriver = "overlay"

func createSecrets(
	ctx context.Context,
	dockerCli *docker_cli_command.DockerCli,
	namespace docker_cli_compose_convert.Namespace,
	secrets []docker_api_types_swarm.SecretSpec,
) error {
	client := dockerCli.Client()

	for _, secretSpec := range secrets {
		secret, _, err := client.SecretInspectWithRaw(ctx, secretSpec.Name)
		if err == nil {
			// secret already exists, then we update that
			if err := client.SecretUpdate(ctx, secret.ID, secret.Meta.Version, secretSpec); err != nil {
				return err
			}
		} else if docker_client.IsErrSecretNotFound(err) {
			// secret does not exist, then we create a new one.
			if _, err := client.SecretCreate(ctx, secretSpec); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func createNetworks(
	ctx context.Context,
	dockerCli *docker_cli_command.DockerCli,
	namespace docker_cli_compose_convert.Namespace,
	networks map[string]docker_api_types.NetworkCreate,
) error {
	client := dockerCli.Client()

	existingNetworks, err := getStackNetworks(ctx, client, namespace.Name())
	if err != nil {
		return err
	}

	existingNetworkMap := make(map[string]docker_api_types.NetworkResource)
	for _, network := range existingNetworks {
		existingNetworkMap[network.Name] = network
	}

	for internalName, createOpts := range networks {
		name := namespace.Scope(internalName)
		if _, exists := existingNetworkMap[name]; exists {
			continue
		}

		if createOpts.Driver == "" {
			createOpts.Driver = defaultNetworkDriver
		}

		fmt.Fprintf(dockerCli.Out(), "Creating network %s\n", name)
		if _, err := client.NetworkCreate(ctx, name, createOpts); err != nil {
			return err
		}
	}

	return nil
}

func deployServices(
	ctx context.Context,
	dockerCli *docker_cli_command.DockerCli,
	services map[string]docker_api_types_swarm.ServiceSpec,
	namespace docker_cli_compose_convert.Namespace,
	sendAuth bool,
) error {
	apiClient := dockerCli.Client()
	out := dockerCli.Out()

	existingServices, err := getStackServices(ctx, apiClient, namespace.Name())
	if err != nil {
		return err
	}

	existingServiceMap := make(map[string]docker_api_types_swarm.Service)
	for _, service := range existingServices {
		existingServiceMap[service.Spec.Name] = service
	}

	for internalName, serviceSpec := range services {
		name := namespace.Scope(internalName)

		encodedAuth := ""
		if sendAuth {
			// Retrieve encoded auth token from the image reference
			image := serviceSpec.TaskTemplate.ContainerSpec.Image
			encodedAuth, err = docker_cli_command.RetrieveAuthTokenFromImage(ctx, dockerCli, image)
			if err != nil {
				return err
			}
		}

		if service, exists := existingServiceMap[name]; exists {
			fmt.Fprintf(out, "Updating service %s (id: %s)\n", name, service.ID)

			updateOpts := docker_api_types.ServiceUpdateOptions{}
			if sendAuth {
				updateOpts.EncodedRegistryAuth = encodedAuth
			}
			response, err := apiClient.ServiceUpdate(
				ctx,
				service.ID,
				service.Version,
				serviceSpec,
				updateOpts,
			)
			if err != nil {
				return err
			}

			for _, warning := range response.Warnings {
				fmt.Fprintln(dockerCli.Err(), warning)
			}
		} else {
			fmt.Fprintf(out, "Creating service %s\n", name)

			createOpts := docker_api_types.ServiceCreateOptions{}
			if sendAuth {
				createOpts.EncodedRegistryAuth = encodedAuth
			}
			if _, err := apiClient.ServiceCreate(ctx, serviceSpec, createOpts); err != nil {
				return err
			}
		}
	}

	return nil
}
