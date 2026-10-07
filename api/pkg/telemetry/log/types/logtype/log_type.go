package logtype

type logType int8

const (
	unknown logType = iota // invalid
	open_telemetry
)
