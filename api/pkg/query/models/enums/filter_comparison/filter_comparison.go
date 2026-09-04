package filtercomparison

import (
	"fmt"
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

func (f FilterComparison) ValidateValues(values []string) error {
	if !f.RequiresValue() {
		if len(values) > 0 {
			return fmt.Errorf("comparison %s does not accept values", f.String())
		}
		return nil
	}
	if len(values) == 0 {
		return fmt.Errorf("comparison %s requires at least one value", f.String())
	}
	switch f {
	case FilterComparisons.BETWEEN:
		if len(values) != 2 {
			return fmt.Errorf("BETWEEN requires exactly 2 values")
		}
	case FilterComparisons.IN:
	default:
		if len(values) != 1 {
			return fmt.Errorf("comparison %s requires exactly 1 value", f.String())
		}
	}
	return nil
}
