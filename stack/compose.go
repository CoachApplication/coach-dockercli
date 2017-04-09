package stack

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	docker_cli_command "github.com/docker/docker/cli/command"
	docker_cli_compose_convert "github.com/docker/docker/cli/compose/convert"
	docker_cli_compose_loader "github.com/docker/docker/cli/compose/loader"
	docker_cli_compose_types "github.com/docker/docker/cli/compose/types"

	coach_config "github.com/CoachApplication/config"
)

func deployComposeFromCoachConfig(ctx context.Context, filename string, config coach_config.Config, env map[string]string, workingDir string, opts deployOptions) (*docker_cli_compose_types.Config, error) {
	var configMap map[string]interface{}
	res := config.Get(&configMap)

	select {
	case <-res.Finished():
		if !res.Success() {
			if errs := res.Errors(); len(errs) > 0 {
				return nil, errs[len(errs)-1]
			} else {
				return nil, fmt.Errorf("Unknown error occured retrieving compose details from compose source Config")
			}
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	confFiles := []docker_cli_compose_types.ConfigFile{
		docker_cli_compose_types.ConfigFile{
			Filename: filename,
			Config:   configMap,
		},
	}

	comConfDetails := docker_cli_compose_types.ConfigDetails{
		WorkingDir:  "path",
		ConfigFiles: confFiles,
		Environment: env,
	}

	comConf, err := docker_cli_compose_loader.Load(comConfDetails)
	return comConf, err
}

func deployComposeDefault(ctx context.Context, dockerCli *docker_cli_command.DockerCli, opts deployOptions) (*docker_cli_compose_types.Config, error) {
	configDetails, err := getConfigDetails(opts)
	if err != nil {
		return nil, err
	}

	config, err := docker_cli_compose_loader.Load(configDetails)
	if err != nil {
		if fpe, ok := err.(*docker_cli_compose_loader.ForbiddenPropertiesError); ok {
			return nil, fmt.Errorf("Compose file contains unsupported options:\n\n%s\n",
				propertyWarnings(fpe.Properties))
		}

		return nil, err
	}

	unsupportedProperties := docker_cli_compose_loader.GetUnsupportedProperties(configDetails)
	if len(unsupportedProperties) > 0 {
		fmt.Fprintf(dockerCli.Err(), "Ignoring unsupported options: %s\n\n",
			strings.Join(unsupportedProperties, ", "))
	}

	deprecatedProperties := docker_cli_compose_loader.GetDeprecatedProperties(configDetails)
	if len(deprecatedProperties) > 0 {
		fmt.Fprintf(dockerCli.Err(), "Ignoring deprecated options:\n\n%s\n\n",
			propertyWarnings(deprecatedProperties))
	}

	if err := checkDaemonIsSwarmManager(ctx, dockerCli); err != nil {
		return nil, err
	}

	return config, nil
}

func getConfigDetails(opts deployOptions) (docker_cli_compose_types.ConfigDetails, error) {
	var details docker_cli_compose_types.ConfigDetails
	var err error

	details.WorkingDir, err = os.Getwd()
	if err != nil {
		return details, err
	}

	configFile, err := getConfigFile(opts.composefile)
	if err != nil {
		return details, err
	}
	// TODO: support multiple files
	details.ConfigFiles = []docker_cli_compose_types.ConfigFile{*configFile}
	details.Environment, err = buildEnvironment(os.Environ())
	if err != nil {
		return details, err
	}
	return details, nil
}

func buildEnvironment(env []string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for _, s := range env {
		// if value is empty, s is like "K=", not "K".
		if !strings.Contains(s, "=") {
			return result, fmt.Errorf("unexpected environment %q", s)
		}
		kv := strings.SplitN(s, "=", 2)
		result[kv[0]] = kv[1]
	}
	return result, nil
}

func getConfigFile(filename string) (*docker_cli_compose_types.ConfigFile, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := docker_cli_compose_loader.ParseYAML(bytes)
	if err != nil {
		return nil, err
	}
	return &docker_cli_compose_types.ConfigFile{
		Filename: filename,
		Config:   config,
	}, nil
}

/**
 * Actual deploy
 */

func deployComposeConfig(ctx context.Context, dockerCli *docker_cli_command.DockerCli, config *docker_cli_compose_types.Config, opts deployOptions) error {
	namespace := docker_cli_compose_convert.NewNamespace(opts.namespace)

	if opts.prune {
		services := map[string]struct{}{}
		for _, service := range config.Services {
			services[service.Name] = struct{}{}
		}
		pruneServices(ctx, dockerCli, namespace, services)
	}

	serviceNetworks := getServicesDeclaredNetworks(config.Services)
	networks, externalNetworks := docker_cli_compose_convert.Networks(namespace, config.Networks, serviceNetworks)
	if err := validateExternalNetworks(ctx, dockerCli, externalNetworks); err != nil {
		return err
	}
	if err := createNetworks(ctx, dockerCli, namespace, networks); err != nil {
		return err
	}

	secrets, err := docker_cli_compose_convert.Secrets(namespace, config.Secrets)
	if err != nil {
		return err
	}
	if err := createSecrets(ctx, dockerCli, namespace, secrets); err != nil {
		return err
	}

	services, err := docker_cli_compose_convert.Services(namespace, config, dockerCli.Client())
	if err != nil {
		return err
	}
	return deployServices(ctx, dockerCli, services, namespace, opts.sendRegistryAuth)
}

func getServicesDeclaredNetworks(serviceConfigs []docker_cli_compose_types.ServiceConfig) map[string]struct{} {
	serviceNetworks := map[string]struct{}{}
	for _, serviceConfig := range serviceConfigs {
		if len(serviceConfig.Networks) == 0 {
			serviceNetworks["default"] = struct{}{}
			continue
		}
		for network := range serviceConfig.Networks {
			serviceNetworks[network] = struct{}{}
		}
	}
	return serviceNetworks
}

func propertyWarnings(properties map[string]string) string {
	var msgs []string
	for name, description := range properties {
		msgs = append(msgs, fmt.Sprintf("%s: %s", name, description))
	}
	sort.Strings(msgs)
	return strings.Join(msgs, "\n\n")
}
