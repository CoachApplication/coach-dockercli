package configwrapper_test

import (
	"testing"

	"github.com/CoachApplication/config/provider/buffered"
	"github.com/CoachApplication/config/provider/yaml"
	dcli_cw "github.com/CoachApplication/handler-dockercli/configwrapper"
	"context"
	"time"
	"golang.org/x/tools/go/gcimporter15/testdata"
)

var cbTest []byte = []byte(`settings: {}

commands:
  hello-world:
    ui:
      label: Hello World
      description: Run the docker hello-world test
      help:8-
        The hello-world test command.

        This test command just runs the hello world image, and
        should capture the hello world information to the output
        passed into it.
    container:
      image: hello-world
      privileged: false
      volumes: []
      links: []
`)

var cTest = dcli_cw.Config{
	Settings: dcli_cw.Settings{},
	Commands: dcli_cw.Commands{
		Commands: map[string]dcli_cw.Command{
			"hello-world": dcli_cw.Command{
				Ui: dcli_cw.CommandUi{
					Label:       "Hello World",
					Description: "Run the docker hello-world test",
					Help: `The hello-world test command.

This test command just runs the hello world image, and
should capture the hello world information to the output
passed into it.`,
				},
				Container: dcli_cw.CommandContainer{
					Image:      "hello-world",
					Privileged: false,
					Volumes:    []string{},
					Links:      []string{},
				},
			},
		},
	},
}

func TestConfig_FromYaml(t *testing.T) {
	dur, _ := time.ParseDuration("2s")
	c := yaml.NewConfig("key", "scope", buffered.NewSingle("key", "scope", cbTest)).Config()
	var cs dcli_cw.Config

	res := c.Get(&cs)

	ctx, _ := context.WithTimeout(context.Background(), dur)
	select {
	case <-res.Finished():

		if !res.Success() {
			t.Error("ConfigProject marshalling using Yaml config failed: ", res.Errors())
		} else if cs.Commands != "test" {
			t.Error("ConfigProject marshalling using Yaml gave the wrong Name: ", p.Name())
		}

	case <-ctx.Done():
		t.Error("Config Get timed out: ", ctx.Err().Error())
	}
}

func TestConfig_ToYaml(t *testing.T) {

}

func TestConfig_Get(t *testing.T) {

}

func TestConfig_Order(t *testing.T) {

}
