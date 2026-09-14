package sqlclientsslmode

type (
	sqlClientSSLMode int8
)

const (
	unknown sqlClientSSLMode = iota // invalid
	disabled
	require
	verifyCA
	verifyFull
)
