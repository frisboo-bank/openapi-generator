package contracts

type Configurable interface {
	GetEnabled() bool
	GetLogger() string
	SetDefaults()
	Validate() error
}
