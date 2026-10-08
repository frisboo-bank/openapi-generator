//go:generate goenums -f -c ./models/enums/log_type/log_type.go
package logtype

type logType int8

const (
	unknown logType = iota // invalid
	open_telemetry
)
