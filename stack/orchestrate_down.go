package stack

import (
	api "github.com/CoachApplication/api"
	base "github.com/CoachApplication/base"
	handler_dockercli "github.com/CoachApplication/handler-dockercli"
)

const (
	OPERATION_ID_ORCHESTRATE_DOWN = "orchestrate.down"
)

type OrchestrateDownOperation struct {
	handler_dockercli.ClientOperationBase
}

func NewOrchestrateDownOperation(base handler_dockercli.ClientOperationBase) *OrchestrateDownOperation {
	return &OrchestrateDownOperation{
		ClientOperationBase: base,
	}
}

func (odo *OrchestrateDownOperation) Operation() api.Operation {
	return api.Operation(odo)
}

func (odo *OrchestrateDownOperation) Id() string {
	return OPERATION_ID_ORCHESTRATE_DOWN
}

func (odo *OrchestrateDownOperation) Ui() api.Ui {
	return base.NewUi(
		odo.Id(),
		"Orchestrate up",
		"Bring up the application app stack",
		"",
	)
}

func (odo *OrchestrateDownOperation) Usage() api.Usage {
	return (&base.ExternalOperationUsage{}).Usage()
}

func (odo *OrchestrateDownOperation) Properties() api.Properties {
	props := base.NewProperties()

	return props.Properties()
}

func (odo *OrchestrateDownOperation) Validate(props api.Properties) api.Result {
	return base.MakeSuccessfulResult()
}

func (odo *OrchestrateDownOperation) Exec(props api.Properties) api.Result {
	return base.MakeSuccessfulResult()
}
