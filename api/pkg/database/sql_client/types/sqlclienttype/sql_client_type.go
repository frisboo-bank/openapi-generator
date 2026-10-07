package sqlclienttype

type sqlClientType int8

const (
	unknown sqlClientType = iota // invalid
	postgresx
	sqlite3x
)
