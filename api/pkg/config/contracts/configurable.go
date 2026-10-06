package contracts

type Configurable interface {
	GetLogger() string
	SetDefaults()
	Validate() error
}
