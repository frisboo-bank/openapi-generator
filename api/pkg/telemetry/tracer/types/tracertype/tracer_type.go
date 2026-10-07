package tracertype

type tracerType int8

const (
	unknown tracerType = iota // invalid
	noop
	open_telemetry
)
