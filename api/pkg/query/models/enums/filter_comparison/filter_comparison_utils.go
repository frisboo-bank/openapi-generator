package filtercomparison

import "fmt"

type FilterOptions struct {
	RequiresValue bool
	MinValues     int
	MaxValues     int
	FormatHint    string
}

func (f FilterComparison) Options() FilterOptions {
	switch f {
	case FilterComparisons.BETWEEN:
		return FilterOptions{
			RequiresValue: true,
			MinValues:     2,
			MaxValues:     2,
			FormatHint:    "two values: lower and upper bound (e.g., '10', '20')",
		}

	case FilterComparisons.IN:
		return FilterOptions{
			RequiresValue: true,
			MinValues:     1,
			MaxValues:     5,
			FormatHint:    "one or more values (e.g., 'apple', 'orange')",
		}

	case FilterComparisons.EQUAL,
		FilterComparisons.GREATER,
		FilterComparisons.GREATEROREQUAL,
		FilterComparisons.LESS,
		FilterComparisons.LESSOREQUAL,
		FilterComparisons.NOTEQUAL:
		return FilterOptions{
			RequiresValue: true,
			MinValues:     1,
			MaxValues:     1,
			FormatHint:    "single value required",
		}

	default:
		return FilterOptions{
			RequiresValue: false,
			FormatHint:    "no value allowed",
		}
	}
}

func (f FilterComparison) ValidateValues(values []string) error {
	options := f.Options()

	if !options.RequiresValue {
		if len(values) > 0 {
			return fmt.Errorf("comparison %q: %s", f.String(), options.FormatHint)
		}
		return nil
	}

	if len(values) < options.MinValues {
		return fmt.Errorf("filter %q: expects at least %d value(s); %s (got %d)",
			f.String(),
			options.MinValues,
			options.FormatHint,
			len(values))
	}

	if options.MaxValues > 0 && len(values) > options.MaxValues {
		return fmt.Errorf("filter %q: expects at most %d value(s); %s (got %d)",
			f.String(),
			options.MaxValues,
			options.FormatHint,
			len(values))
	}

	return nil
}
