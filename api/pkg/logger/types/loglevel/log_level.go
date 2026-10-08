//go:generate goenums -f -c ./log_level.go
package loglevel

type logLevel int8

const (
	unknown logLevel = iota // invalid
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	panicLevel
	fatalLevel
	traceLevel
)
