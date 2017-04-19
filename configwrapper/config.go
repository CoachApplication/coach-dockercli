package configwrapper

import "github.com/CoachApplication/command"

type Config struct {
	Settings ConfigSettings `yaml:"settings"`
	Commands map[string]ConfigCommand `yaml:"commands"`
}

func (c *Config) Get(id string) (command.Command, error) {
	
}
func (c *Config) Order() []string {

}

type ConfigSettings struct {
	ListPermissionApply bool `yaml:"OnlyShowPermitted"`
}


type ConfigCommandList struct {
	commands map[string]ConfigCommand
}

func (ccp *ConfigCommandList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var Coms map[string]ConfigCommand

	if err := unmarshal(&Coms); err != nil {
		return err
	}
	for id, com := range Coms {
		ccp.commands[id] = com
	}
}

type ConfigCommand struct {
	Id string
	Privileged bool `yaml:"privileged"`
	Image string `yaml:"string"`
	Volumes []string `yaml:"volumes"`
	Links []string `yaml:"links"`
}
