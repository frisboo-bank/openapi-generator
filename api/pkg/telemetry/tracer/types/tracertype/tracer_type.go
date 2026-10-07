package tracertype

type tracerType int8

const (
	unknown tracerType = iota // invalid
	open_telemetry
)
