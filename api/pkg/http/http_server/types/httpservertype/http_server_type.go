//go:generate goenums -f -c ./http_server_type.go
package httpservertype

type httpServerType int8

const (
	unknown httpServerType = iota // invalid
	echo
)
