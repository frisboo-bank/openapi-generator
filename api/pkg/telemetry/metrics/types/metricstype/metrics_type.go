//go:generate goenums -f -c ./metrics_type.go
package metricstype

type metricsType int8

const (
	unknown metricsType = iota // invalid
	open_telemetry
)
