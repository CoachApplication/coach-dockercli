package configwrapper

import (
	"github.com/CoachApplication/command"
	"github.com/CoachApplication/command/provider"
	"github.com/CoachApplication/config"
	dockerhandler_command "github.com/CoachApplication/handler-dockercli/command"
)

const (
	CONFIG_ID_COMMANDS = "commands"
)

type Provider struct {
	factory dockerhandler_command.Factory
	Config  Config
}

func NewProviderFromConfigWrapper(wr config.Wrapper, factory dockerhandler_command.Factory) *Provider {

	conf := NewConfig()
	if scopedConfs, err := wr.Get(CONFIG_ID_COMMANDS); err == nil {
		for _, scope := range scopedConfs.Order() {
			scopeConf, _ := scopedConfs.Get(scope)

			res := scopeConf.HasValue()
			<-res.Finished() // @TODO compete with a context to avoid breaking on config timeout

			if res.Success() {
				var scopeCommandConf Config
				scopeConf.Get(&scopeCommandConf)

				conf.Merge(scopeCommandConf)
			}
		}
	}

	return
}

func (p *Provider) Provider() provider.Provider {
	return provider.Provider(p)
}

func (p *Provider) Get(id string) (command.Command, error) {

}
func (p *Provider) Order() []string {

}
