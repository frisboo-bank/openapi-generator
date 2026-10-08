//go:generate goenums -f -c ./tracer_type.go
package tracertype

type tracerType int8

const (
	unknown tracerType = iota // invalid
	open_telemetry
)
