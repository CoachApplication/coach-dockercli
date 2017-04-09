package stack

import (
	"context"
	"errors"
	"fmt"

	docker_api_types "github.com/docker/docker/api/types"
	docker_api_types_filters "github.com/docker/docker/api/types/filters"
	docker_api_types_swarm "github.com/docker/docker/api/types/swarm"
	docker_cli_command "github.com/docker/docker/cli/command"
	docker_cli_compose_convert "github.com/docker/docker/cli/compose/convert"
	docker_client "github.com/docker/docker/client"
	docker_opts "github.com/docker/docker/opts"
)

// checkDaemonIsSwarmManager does an Info API call to verify that the daemon is
// a swarm manager. This is necessary because we must create networks before we
// create services, but the API call for creating a network does not return a
// proper status code when it can't create a network in the "global" scope.
func checkDaemonIsSwarmManager(ctx context.Context, dockerCli *docker_cli_command.DockerCli) error {
	info, err := dockerCli.Client().Info(ctx)
	if err != nil {
		return err
	}
	if !info.Swarm.ControlAvailable {
		return errors.New("This node is not a swarm manager. Use \"docker swarm init\" or \"docker swarm join\" to connect this node to swarm and try again.")
	}
	return nil
}

func getStackFilter(namespace string) docker_api_types_filters.Args {
	filter := docker_api_types_filters.NewArgs()
	filter.Add("label", docker_cli_compose_convert.LabelNamespace+"="+namespace)
	return filter
}

func getStackFilterFromOpt(namespace string, opt docker_opts.FilterOpt) docker_api_types_filters.Args {
	filter := opt.Value()
	filter.Add("label", docker_cli_compose_convert.LabelNamespace+"="+namespace)
	return filter
}

func getAllStacksFilter() docker_api_types_filters.Args {
	filter := docker_api_types_filters.NewArgs()
	filter.Add("label", docker_cli_compose_convert.LabelNamespace)
	return filter
}

func getStackServices(
	ctx context.Context,
	apiclient docker_client.APIClient,
	namespace string,
) ([]docker_api_types_swarm.Service, error) {
	return apiclient.ServiceList(
		ctx,
		docker_api_types.ServiceListOptions{Filters: getStackFilter(namespace)})
}

func getStackNetworks(
	ctx context.Context,
	apiclient docker_client.APIClient,
	namespace string,
) ([]docker_api_types.NetworkResource, error) {
	return apiclient.NetworkList(
		ctx,
		docker_api_types.NetworkListOptions{Filters: getStackFilter(namespace)})
}

func getStackSecrets(
	ctx context.Context,
	apiclient docker_client.APIClient,
	namespace string,
) ([]docker_api_types_swarm.Secret, error) {
	return apiclient.SecretList(
		ctx,
		docker_api_types.SecretListOptions{Filters: getStackFilter(namespace)})
}

// pruneServices removes services that are no longer referenced in the source
func pruneServices(ctx context.Context, dockerCli docker_cli_command.Cli, namespace docker_cli_compose_convert.Namespace, services map[string]struct{}) bool {
	client := dockerCli.Client()

	oldServices, err := getStackServices(ctx, client, namespace.Name())
	if err != nil {
		fmt.Fprintf(dockerCli.Err(), "Failed to list services: %s", err)
		return true
	}

	pruneServices := []docker_api_types_swarm.Service{}
	for _, service := range oldServices {
		if _, exists := services[namespace.Descope(service.Spec.Name)]; !exists {
			pruneServices = append(pruneServices, service)
		}
	}
	return removeServices(ctx, dockerCli, pruneServices)
}

func validateExternalNetworks(
	ctx context.Context,
	dockerCli *docker_cli_command.DockerCli,
	externalNetworks []string) error {
	client := dockerCli.Client()

	for _, networkName := range externalNetworks {
		network, err := client.NetworkInspect(ctx, networkName, false)
		if err != nil {
			if docker_client.IsErrNetworkNotFound(err) {
				return fmt.Errorf("network %q is declared as external, but could not be found. You need to create the network before the stack is deployed (with overlay driver)", networkName)
			}
			return err
		}
		if network.Scope != "swarm" {
			return fmt.Errorf("network %q is declared as external, but it is not in the right scope: %q instead of %q", networkName, network.Scope, "swarm")
		}
	}

	return nil
}
