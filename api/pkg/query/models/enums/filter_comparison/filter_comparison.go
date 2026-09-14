package filtercomparison

type (
	filterComparison int8
)

const (
	unknown filterComparison = iota // invalid
	// with value
	between
	equal
	greater
	greaterOrEqual
	in
	less
	lessOrEqual
	notEqual

	// no value
	isEmpty
	isFalse
	isNotEmpty
	isNotFalse
	isNotNull
	isNotTrue
	isNull
	isTrue
	isUnknown
)
