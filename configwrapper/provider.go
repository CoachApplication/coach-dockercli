package configwrapper

import (
	"errors"

	"github.com/CoachApplication/command"
	"github.com/CoachApplication/command/provider"
	"github.com/CoachApplication/config"
	dockerhandler_command "github.com/CoachApplication/handler-dockercli/command"
)

const (
	CONFIG_ID_COMMANDS = "commands"
)

// Provider A docker cli command prpovider that uses docker to run commands in a container
type Provider struct {
	factory dockerhandler_command.Factory // must be provided
	configWrapper config.Wrapper // must be provided

	config  Config // generated when needed
}

func NewProviderFromConfigWrapper(wr config.Wrapper, fac *dockerhandler_command.Factory) *Provider {
	if fac == nil {
		fac = dockerhandler_command.NewDefaultFactory()
	}

	return &Provider{
		configWrapper: wr,
		factory: *fac,
	}
}

func (p *Provider) Provider() provider.Provider {
	return provider.Provider(p)
}

func (p *Provider) safe() {
	if &p.config == nil {
		p.load()
	}
}

func (p *Provider) load() {
	wr := p.configWrapper
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

	p.config = *conf
}

func (p *Provider) Get(id string) (command.Command, error) {
	p.safe()
	if confCom, err := p.config.CommandGet(id); err != nil {
		return nil, errors.New("Command not found: "+id)
	} else {
		return confCom.CommandFromFactory(p.factory)
	}
}
func (p *Provider) Order() []string {
	return p.config.CommandOrder()
}
