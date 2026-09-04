package rpcservertype

type (
	rpcServerType int8
)

const (
	unknown rpcServerType = iota // invalid
	grpc
)
