package cachetype

type (
	cacheType int8
)

const (
	unknown cacheType = iota // invalid
	redis
	memory
)
