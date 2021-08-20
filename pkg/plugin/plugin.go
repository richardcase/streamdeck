package plugin

type Plugin interface {
	Configure() error
	GetActions() ([]Action, error)
}

type Action interface {
	ID() string
	Description() string
	Pressed()
}
