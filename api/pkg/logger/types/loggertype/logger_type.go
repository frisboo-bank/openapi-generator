//go:generate goenums -f -c ./logger_type.go
package loggertype

type loggerType int8

const (
	unknown loggerType = iota
	zerolog
)
