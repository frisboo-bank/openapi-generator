package environment

type (
	environment int8
)

const (
	unknown environment = iota // invalid
	development
	preprod
	production
	testing
)

func (e Environment) IsEnvironment(env Environment) bool {
	return e == env
}

func (e Environment) IsDevelopment() bool {
	return e == Environments.DEVELOPMENT
}

func (e Environment) IsTesting() bool {
	return e == Environments.TESTING
}

func (e Environment) IsPreprod() bool {
	return e == Environments.PREPROD
}

func (e Environment) IsProduction() bool {
	return e == Environments.PRODUCTION
}
