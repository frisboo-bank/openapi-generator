//go:generate goenums -f -c ./hook_status.go
package hookstatus

type hookStatus int8

const (
	unknown hookStatus = iota
	idle
	running
	stopped
)
