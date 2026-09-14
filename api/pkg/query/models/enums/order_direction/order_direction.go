package orderdirection

type (
	orderDirection int8
)

const (
	unknown orderDirection = iota // invalid
	asc
	desc
)
