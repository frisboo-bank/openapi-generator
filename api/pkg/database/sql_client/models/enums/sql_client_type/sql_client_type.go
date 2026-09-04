package sqlclienttype

type (
	sqlClientType int8
)

const (
	unknown sqlClientType = iota // invalid
	postgres
	postgresx
	sqlite3x
)
