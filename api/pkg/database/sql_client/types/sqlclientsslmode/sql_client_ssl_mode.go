//go:generate goenums -f -c ./sql_client_ssl_mode.go
package sqlclientsslmode

type sqlClientSSLMode int8

const (
	unknown sqlClientSSLMode = iota // invalid
	disabled
	require
	verifyCA
	verifyFull
)
