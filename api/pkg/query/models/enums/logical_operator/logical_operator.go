package logicaloperator

type (
	logicalOperator int8
)

const (
	unknown logicalOperator = iota // invalid
	and
	or
)
