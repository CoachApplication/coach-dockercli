package configwrapper

import (
	"errors"
	handler_docker_command "github.com/CoachApplication/handler-dockercli/command"
	"github.com/CoachApplication/command"
	"github.com/CoachApplication/base"
)

type Config struct {
	Settings Settings `yaml:"settings"`
	Commands Commands `yaml:"commands"`
}

func NewConfig() *Config {
	return &Config{
		Settings: Settings{},
		Commands: Commands{},
	}
}

func (c *Config) Merge(merge Config) {
	c.Settings.Merge(merge.Settings)
	c.Commands.Merge(merge.Commands)
}

func (c *Config) CommandGet(id string) (Command, error) {
	return c.Commands.Get(id)
}
func (c *Config) CommandOrder() []string {
	return c.Commands.Order()
}

type Settings struct {
	ListPermissionApply bool `yaml:"OnlyShowPermitted"`
}

func (s *Settings) Merge(merge Settings) {

}

type Commands struct {
	Commands map[string]Command
}

func (c *Commands) Merge(merge Commands) {

}

func (cs *Commands) Get(id string) (Command, error) {
	if com, found := cs.Commands[id]; found {
		return com, nil
	} else {
		return Command{}, errors.New("Command not found")
	}

}
func (cs *Commands) Order() []string {
	return []string{}
}

func (cs *Commands) AddCommand(id string, cmd Command) {
	cs.Commands[id] = cmd
}

func (cs *Commands) AddCommandComponents(id string, cmdui CommandUi, cmdCont CommandContainer) {
	cs.Commands[id] = Command{
		Ui:        cmdui,
		Container: cmdCont,
	}
}

func (cs *Commands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var Coms map[string]Command

	if err := unmarshal(&Coms); err != nil {
		return err
	}
	for id, com := range Coms {
		cs.AddCommand(id, com)
	}
	return nil
}

type Command struct {
	Ui        CommandUi        `yaml:"ui"`
	Container CommandContainer `yaml:"container"`
}

// Use a Command factory to convert this Config Command to a core command.Command
func (c *Command) CommandFromFactory(f handler_docker_command.Factory) (command.Command, error) {
	return f.NewCommand(
		c.Container.Id,
		base.NewUi(
			c.Container.Id,
			c.Ui.Label,
			c.Ui.Description,
			c.Ui.Help,
		),
		base.ExternalOperationUsage{}.Usage(), // For now just assume that all commands are public
		c.Container.Privileged,
		c.Container.Image,
		c.Container.Volumes,
		c.Container.Links,
	), nil
}

type CommandUi struct {
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Help        string `yaml:"help"`
}

type CommandContainer struct {
	Id         string
	Privileged bool     `yaml:"privileged"`
	Image      string   `yaml:"string"`
	Volumes    []string `yaml:"volumes"`
	Links      []string `yaml:"links"`
}
