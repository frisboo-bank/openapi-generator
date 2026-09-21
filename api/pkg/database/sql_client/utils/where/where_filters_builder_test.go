package where_test

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/where"
	"frisboo-bank/openapi-generator-service/pkg/query"
	filtercomparisonenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNamedWhereFiltersClause(t *testing.T) {
	examples := []struct {
		name    string
		where   string
		filters []query.Filter
		res     string
		resArgs map[string]any
		isErr   bool
	}{
		{
			name:    "empty filters",
			where:   "hidden_at IS NULL",
			filters: []query.Filter{},
			res:     "hidden_at IS NULL",
			resArgs: nil,
		},
		{
			name:    "simple filter",
			where:   "",
			filters: []query.Filter{{Field: "username", Values: []string{"john.doe"}, Comparator: filtercomparisonenum.FilterComparisons.EQUAL}},
			res:     "username = :username_0",
			resArgs: map[string]any{"username_0": "john.doe"},
		},
		{
			name:    "append filter to existing where",
			where:   "hidden_at IS NULL",
			filters: []query.Filter{{Field: "username", Values: []string{"john.doe"}, Comparator: filtercomparisonenum.FilterComparisons.EQUAL}},
			res:     "hidden_at IS NULL AND username = :username_0",
			resArgs: map[string]any{"username_0": "john.doe"},
		},
		{
			name:  "multiple filters combined with AND",
			where: "hidden_at IS NULL",
			filters: []query.Filter{
				{Field: "slug", Comparator: filtercomparisonenum.FilterComparisons.EQUAL, Values: []string{"foo"}},
				{Field: "created_at", Comparator: filtercomparisonenum.FilterComparisons.GREATER, Values: []string{"2024-01-01"}},
			},
			res:     "hidden_at IS NULL AND slug = :slug_0 AND created_at > :created_at_1",
			resArgs: map[string]any{"slug_0": "foo", "created_at_1": "2024-01-01"},
		},
		{
			name:  "all comparison operators",
			where: "",
			filters: []query.Filter{
				{Field: "a", Comparator: filtercomparisonenum.FilterComparisons.EQUAL, Values: []string{"1"}},
				{Field: "b", Comparator: filtercomparisonenum.FilterComparisons.NOTEQUAL, Values: []string{"2"}},
				{Field: "c", Comparator: filtercomparisonenum.FilterComparisons.GREATER, Values: []string{"3"}},
				{Field: "d", Comparator: filtercomparisonenum.FilterComparisons.GREATEROREQUAL, Values: []string{"4"}},
				{Field: "e", Comparator: filtercomparisonenum.FilterComparisons.LESS, Values: []string{"5"}},
				{Field: "f", Comparator: filtercomparisonenum.FilterComparisons.LESSOREQUAL, Values: []string{"6"}},
			},
			res: "a = :a_0 AND b != :b_1 AND c > :c_2 AND d >= :d_3 AND e < :e_4 AND f <= :f_5",
			resArgs: map[string]any{
				"a_0": "1", "b_1": "2", "c_2": "3",
				"d_3": "4", "e_4": "5", "f_5": "6",
			},
		},
		{
			name:  "IN with multiple values",
			where: "",
			filters: []query.Filter{
				{Field: "status", Comparator: filtercomparisonenum.FilterComparisons.IN, Values: []string{"active", "pending", "archived"}},
			},
			res:     "status IN (:status_0_0, :status_0_1, :status_0_2)",
			resArgs: map[string]any{"status_0_0": "active", "status_0_1": "pending", "status_0_2": "archived"},
		},
		{
			name:  "BETWEEN with two values",
			where: "",
			filters: []query.Filter{
				{Field: "created_at", Comparator: filtercomparisonenum.FilterComparisons.BETWEEN, Values: []string{"2024-01-01", "2024-12-31"}},
			},
			res:     "created_at BETWEEN :created_at_0_low AND :created_at_0_high",
			resArgs: map[string]any{"created_at_0_low": "2024-01-01", "created_at_0_high": "2024-12-31"},
		},
		{
			name:  "nullary operators without values",
			where: "1=1",
			filters: []query.Filter{
				{Field: "hidden_at", Comparator: filtercomparisonenum.FilterComparisons.ISNULL, Values: nil},
				{Field: "deleted_at", Comparator: filtercomparisonenum.FilterComparisons.ISNOTNULL, Values: []string{}},
			},
			res:     "1=1 AND hidden_at IS NULL AND deleted_at IS NOT NULL",
			resArgs: map[string]any{},
		},

		// validation
		{
			name:  "BETWEEN with wrong value count errors",
			where: "",
			filters: []query.Filter{
				{Field: "age", Comparator: filtercomparisonenum.FilterComparisons.BETWEEN, Values: []string{"18"}},
			},
			res:   "BETWEEN requires exactly 2 values",
			isErr: true,
		},
		{
			name:  "IN with empty values errors",
			where: "",
			filters: []query.Filter{
				{Field: "status", Comparator: filtercomparisonenum.FilterComparisons.IN, Values: []string{}},
			},
			res:   "IN requires at least 1 value",
			isErr: true,
		},
		{
			name:  "unary with wrong value count errors",
			where: "",
			filters: []query.Filter{
				{Field: "slug", Comparator: filtercomparisonenum.FilterComparisons.EQUAL, Values: []string{"a", "b"}},
			},
			res:   "requires exactly 1 value",
			isErr: true,
		},
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			res, args, err := where.BuildNamedWhereFiltersClause(e.where, e.filters)

			if e.isErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), e.res)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, e.res, res)
			assert.Equal(t, e.resArgs, args)
		})
	}
}
