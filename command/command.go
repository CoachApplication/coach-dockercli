package command

import (
	"context"
	"io"
	"math/rand"
	"os"
	"time"

	docker_api_types "github.com/docker/docker/api/types"
	docker_api_types_network "github.com/docker/docker/api/types/network"
	docker_client "github.com/docker/docker/client"

	"errors"
	"github.com/CoachApplication/api"
	"github.com/CoachApplication/base"
	"github.com/CoachApplication/command"
	"github.com/CoachApplication/handler-dockercli"
	"github.com/docker/docker/api/types/container"
)

// Command a command.Command which will execute a command inside a container
type Command struct {
	id    string
	usage api.Usage
	ui    api.Ui

	ctx context.Context

	cli           docker_client.APIClient
	networkDriver string

	config     container.Config
	hostConfig container.HostConfig

	out io.WriteCloser
	err io.WriteCloser
	in  io.ReadCloser
}

func NewCommand(
	id string,
	ui api.Ui,
	usage api.Usage,
	ctx context.Context,
	cli docker_client.APIClient,
	networkDriver string,
	config container.Config,
	hostConfig container.HostConfig,
	out,
	err io.WriteCloser,
	in io.ReadCloser,
) *Command {
	if out == nil {
		out = os.Stdout
	}
	if err == nil {
		err = os.Stderr
	}
	if in == nil {
		in = os.Stdin
	}
	if cli == nil {
		cli, _ = dockercli.DefaultClient()
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if networkDriver == "" {
		networkDriver = "overlay"
	}

	return &Command{
		id:    id,
		ui:    ui,
		usage: usage,

		ctx: ctx,

		cli:           cli,
		networkDriver: networkDriver,

		config:     config,
		hostConfig: hostConfig,

		in:  in,
		out: out,
		err: err,
	}
}

func (c *Command) Command() command.Command {
	return command.Command(c)
}

// Id Unique string machine name identifier for the Operation
func (c *Command) Id() string {
	return c.id
}

// UI Return a UI interaction definition for the Operation
func (c *Command) Ui() api.Ui {
	return c.ui
}

// Usage Define how the Operation is intended to be executed.
func (c *Command) Usage() api.Usage {
	return c.usage
}

// Properties provide the expected Operation with default values
func (c *Command) Properties() api.Properties {
	props := base.NewProperties()

	props.Add((&command.ArgsProperty{}).Property())

	return props.Properties()
}

// Validate Validate that the Operation can Execute if passed proper Property data
func (c *Command) Validate(props api.Properties) api.Result {
	res := base.NewResult()

	go func(props api.Properties) {
		if _, err := props.Get(command.PROPERTY_ID_COMMAND_ARGS); err != nil {
			res.AddError(err)
			res.MarkFinished()
		}

		res.MarkFinished()
	}(props)

	return res.Result()
}

/**
 * Exec runs the operation from a Properties set, and return a result
 */
func (c *Command) Exec(props api.Properties) api.Result {
	res := base.NewResult()

	go func(props api.Properties) {
		args := []string{}
		if argsProp, err := props.Get(command.PROPERTY_ID_COMMAND_ARGS); err != nil {
			res.AddError(err)
			res.MarkFailed()
		} else {
			args = argsProp.Get().([]string)
		}

		runId := randomIdentifier()

		c.createCommandNetwork(res, runId, map[string]string{})
		defer c.removeCommandNetwork(res, runId)

		c.runCommand(res, runId, args)

		res.MarkFinished()
	}(props)

	return res.Result()
}

func (c *Command) runCommand(res *base.Result, runId string, args []string) {
	networkingConfig := docker_api_types_network.NetworkingConfig{
		EndpointsConfig: map[string]*docker_api_types_network.EndpointSettings{
			"command": &docker_api_types_network.EndpointSettings{
				Aliases:   []string{c.Id()},
				NetworkID: c.networkName(runId),
			},
		},
	}
	attachOptions := docker_api_types.ContainerAttachOptions{
		Stream:     true,
		Stdin:      c.in != nil,
		Stdout:     c.out != nil,
		Stderr:     c.err != nil,
		DetachKeys: "",
		Logs:       true,
	}

	if body, err := c.cli.ContainerCreate(c.ctx, &c.config, &c.hostConfig, &networkingConfig, c.containerName(runId)); err != nil {
		res.AddError(err)
		res.MarkFailed()
	} else if id := body.ID; id == "" {
		for _, warning := range body.Warnings {
			res.AddError(errors.New(warning))
		}
		res.MarkFailed()
	} else if response, err := c.cli.ContainerAttach(c.ctx, id, attachOptions); err != nil {

	} else {
		//(types.HijackedResponse, error)
		defer response.Close()

	}

	res.MarkSucceeded()
}

func (c *Command) containerName(runId string) string {
	return c.Id() + "_" + runId
}
func (c *Command) networkName(runId string) string {
	return c.Id() + "_" + runId
}

func (c *Command) createCommandNetwork(res *base.Result, runId string, links map[string]string) {
	netName := c.networkName(runId)

	createOpts := docker_api_types.NetworkCreate{
		// CheckDuplicate: true,
		Driver: c.networkDriver,
		// EnableIPv6: true,
		Internal:   false, // many commands will need internet access
		Attachable: true,
		Labels: map[string]string{
			"coach.command.network": "yes",
			"coach.command.runId":   runId,
		},
	}

	if _, err := c.cli.NetworkCreate(c.ctx, netName, createOpts); err != nil {
		res.AddError(err)
	} else {

		for id, alias := range links {
			set := docker_api_types_network.EndpointSettings{
				Aliases: []string{alias},
			}

			if err := c.cli.NetworkConnect(c.ctx, netName, id, &set); err != nil {
				res.AddError(err)
			}
		}
	}
}
func (c *Command) removeCommandNetwork(res *base.Result, runId string) {
	if err := c.cli.NetworkRemove(c.ctx, c.networkName(runId)); err != nil {
		res.AddError(err)
	}
}

func randomIdentifier() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 16)
	for i, _ := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
