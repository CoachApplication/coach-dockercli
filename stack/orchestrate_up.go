package stack

import (
	api "github.com/CoachApplication/api"
	base "github.com/CoachApplication/base"
	handler_dockercli "github.com/CoachApplication/handler-dockercli"
)

const (
	OPERATION_ID_ORCHESTRATE_UP = "orchestrate.up"
)

type OrchestrateUpOperation struct {
	handler_dockercli.ClientOperationBase
}

func NewOrchestrateUpOperation(base handler_dockercli.ClientOperationBase) *OrchestrateUpOperation {
	return &OrchestrateUpOperation{
		ClientOperationBase: base,
	}
}

func (ouo *OrchestrateUpOperation) Operation() api.Operation {
	return api.Operation(ouo)
}

func (ouo *OrchestrateUpOperation) Id() string {
	return OPERATION_ID_ORCHESTRATE_UP
}

func (ouo *OrchestrateUpOperation) Ui() api.Ui {
	return base.NewUi(
		ouo.Id(),
		"Orchestrate up",
		"Bring up the application app stack",
		"",
	)
}

func (ouo *OrchestrateUpOperation) Usage() api.Usage {
	return (&base.ExternalOperationUsage{}).Usage()
}

func (ouo *OrchestrateUpOperation) Properties() api.Properties {
	props := base.NewProperties()

	return props.Properties()
}

func (ouo *OrchestrateUpOperation) Validate(props api.Properties) api.Result {
	return base.MakeSuccessfulResult()
}

func (ouo *OrchestrateUpOperation) Exec(props api.Properties) api.Result {
	return base.MakeSuccessfulResult()
}
