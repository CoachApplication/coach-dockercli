package configwrapper

import (
	"github.com/CoachApplication/api"
	"github.com/CoachApplication/base"
	"github.com/CoachApplication/config"

	handler_dockercli "github.com/CoachApplication/handler-dockercli"
	handler_dockercli_stack "github.com/CoachApplication/handler-dockercli/stack"
)

func MakeOrchestrateOperations(wr config.Wrapper) api.Operations {
	ops := base.NewOperations()

	cob := handler_dockercli.NewClientOperationBaseDefault()

	ops.Add(handler_dockercli_stack.NewOrchestrateUpOperation(*cob).Operation())
	ops.Add(handler_dockercli_stack.NewOrchestrateDownOperation(*cob).Operation())

	return ops.Operations()
}
