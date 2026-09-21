package query_test

import (
	"testing"

	"frisboo-bank/openapi-generator-service/pkg/query"
	filtercomparison "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"

	"pgregory.net/rapid"
)

func generateQueryFilter() *rapid.Generator[*query.Filter] {
	return rapid.Custom(func(t *rapid.T) *query.Filter {
		var field string
		for field == "" {
			field = rapid.String().Draw(t, "field")
		}

		comparisons := make([]filtercomparison.FilterComparison, 0)
		for comparison := range filtercomparison.FilterComparisons.All() {
			if comparison.IsValid() {
				comparisons = append(comparisons, comparison)
			}
		}
		comparison := rapid.SampledFrom(comparisons).Draw(t, "comparison")

		return &query.Filter{Field: field, Comparator: comparison}
	})
}

func TestQueryFilters(t *testing.T) {
	t.Run("empty filters is always valid", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			allowed := rapid.SliceOf(rapid.String()).Draw(t, "allowed")

			qo := &query.Query{}

			if err := qo.Validate(allowed); err != nil {
				t.Fatalf("empty filters should validate: %v", err)
			}
		})
	})

	t.Run("filters included not allowed field should fail validation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			filters := rapid.SliceOfN(generateQueryFilter(), 2, 20).Draw(t, "filters")

			badIdx := rapid.IntRange(0, len(filters)-1).Draw(t, "badIdx")
			notAllowed := filters[badIdx].Field

			allowedSet := make(map[string]struct{})
			for i, filter := range filters {
				if i == badIdx {
					continue
				}
				allowedSet[filter.Field] = struct{}{}
			}
			allowed := make([]string, 0, len(allowedSet))
			for field := range allowedSet {
				allowed = append(allowed, field)
			}

			qo := &query.Query{Filters: filters}

			err := qo.Validate(allowed)
			if err == nil {
				t.Fatalf("expected validation error for field %q not in allowed %q", notAllowed, allowed)
			}
		})
	})
}
