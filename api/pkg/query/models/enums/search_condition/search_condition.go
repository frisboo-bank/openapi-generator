package searchcondition

type (
	searchCondition int8
)

const (
	unknown searchCondition = iota // invalid
	contains
	starts
	ends
	equals
)
