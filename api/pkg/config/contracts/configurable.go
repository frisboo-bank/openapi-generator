package contracts

type Configurable interface {
	GetLogger() string
	SetDefaults()
	Validate() error
}

// Enablable is an optional capability: a config carrying an on/off toggle.
type Enablable interface {
	Enable() bool
}
