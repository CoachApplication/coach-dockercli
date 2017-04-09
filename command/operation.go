package command

import (
	"github.com/CoachApplication/command"
	"github.com/CoachApplication/api"
	"github.com/CoachApplication/base"
)

type GetOperation struct {
	command.GetOperation
}


func (gop *GetOperation) Operation() api.Operation {
	return api.Operation(gop)
}

func (gop *GetOperation) Properties() api.Properties {
	props := base.NewProperties()

	props.Add((&command.IdProperty{}).Property())

	return props.Properties()
}

func (gop *GetOperation) Exec(props api.Properties) api.Result {
	res := base.NewResult()


	return res.Result()
}

type ListOperation struct {
	command.ListOperation

	provider Provider
}

func (lo *ListOperation) Operation() api.Operation {
	return api.Operation(lo)
}

func (lo *ListOperation) Properties() api.Properties {
	return base.NewProperties().Properties()
}

func (lo *ListOperation) Exec(props api.Properties) api.Result {
	res := base.NewResult()

	go func(provider Provider) {
		idsProp := command.IdsProperty{}
		idsProp.Set(provider.Order())
		res.AddProperty(idsProp)

		res.MarkSucceeded()
		res.MarkFinished()
	}(lo.provider)

	return res.Result()
}
