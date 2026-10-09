//go:generate goenums -f -c ./rpc_server_type.go
package rpcservertype

type rpcServerType int8

const (
	unknown rpcServerType = iota // invalid
	grpc
)
