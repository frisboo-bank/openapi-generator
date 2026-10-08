//go:generate goenums -f -c ./sql_client_type.go
package sqlclienttype

type sqlClientType int8

const (
	unknown sqlClientType = iota // invalid
	postgresx
	sqlite3x
)
