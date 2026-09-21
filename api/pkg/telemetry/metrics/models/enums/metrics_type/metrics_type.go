package metrictype

type metricsType int8

const (
	unknown metricsType = iota // invalid
	open_telemetry
)
