package metrictype

type metricsType int8

const (
	unknown metricsType = iota // invalid
	noop
	open_telemetry
)
