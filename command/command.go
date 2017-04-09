package command

/**
 * Run commands against a list of services
 */

type Source interface {
	Links() []ServiceContainerMatch
}

type Provider interface {
	Get(id string) (Command, error)
	Order() []string
}
