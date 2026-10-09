//go:generate goenums -f -c ./log_type.go
package logtype

type logType int8

const (
	unknown logType = iota // invalid
	open_telemetry
)
