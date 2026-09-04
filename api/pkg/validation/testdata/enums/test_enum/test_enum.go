package testenum

type (
	testEnum int8
)

const (
	unknown testEnum = iota // invalid
	passed
	skipped
	scheduled
	running
	booked
)
