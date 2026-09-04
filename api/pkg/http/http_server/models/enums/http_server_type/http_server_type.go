package httpservertype

type (
	httpServerType int8
)

const (
	unknown httpServerType = iota // invalid
	echo
)
