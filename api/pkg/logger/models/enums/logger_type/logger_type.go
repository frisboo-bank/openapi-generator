package loggertype

type (
	loggerType int8
)

const (
	unknown loggerType = iota
	noop
	zerolog
)
