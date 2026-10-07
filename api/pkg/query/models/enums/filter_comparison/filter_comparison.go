package filtercomparison

import (
)

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
	lower
	lowerOrEqual
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

func (f FilterComparison) RequiresValue() bool {
	switch f {
	case
		FilterComparisons.ISEMPTY,
		FilterComparisons.ISFALSE,
		FilterComparisons.ISNOTEMPTY,
		FilterComparisons.ISNOTFALSE,
		FilterComparisons.ISNOTNULL,
		FilterComparisons.ISNOTTRUE,
		FilterComparisons.ISNULL,
		FilterComparisons.ISTRUE,
		FilterComparisons.ISUNKNOWN:
		return false
	default:
		return true
	}
}

