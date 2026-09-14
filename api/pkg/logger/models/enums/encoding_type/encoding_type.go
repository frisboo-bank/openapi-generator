package encodingtype

type (
	encodingType int8
)

const (
	unknown encodingType = iota
	json
	text
)
