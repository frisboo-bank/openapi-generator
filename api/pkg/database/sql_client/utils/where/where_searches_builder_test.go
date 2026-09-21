package where_test

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/where"
	"frisboo-bank/openapi-generator-service/pkg/query"
	logicaloperatorenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
	searchconditionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNamedWhereSearchesClause(t *testing.T) {
	examples := []struct {
		name        string
		where       string
		searchGroup *query.SearchGroup
		res         string
		resArgs     map[string]any
		isErr       bool
	}{
		{
			name:        "empty search group",
			where:       "hidden_at IS NULL",
			searchGroup: nil,
			res:         "hidden_at IS NULL",
			resArgs:     nil,
		},
		{
			name:        "empty search group",
			where:       "hidden_at IS NULL",
			searchGroup: &query.SearchGroup{},
			res:         "hidden_at IS NULL",
			resArgs:     nil,
		},
		{
			name:  "single search condition",
			where: "",
			searchGroup: &query.SearchGroup{
				Operator: logicaloperatorenum.LogicalOperators.AND,
				Conditions: []query.Search{
					{Field: "name", Value: "foo", Operator: searchconditionenum.SearchConditions.CONTAINS},
				},
			},
			res:     "name ILIKE :search_d0_name_0",
			resArgs: map[string]any{"search_d0_name_0": "%foo%"},
		},
		{
			name:  "multiple conditions with OR",
			where: "hidden_at IS NULL",
			searchGroup: &query.SearchGroup{
				Operator: logicaloperatorenum.LogicalOperators.OR,
				Conditions: []query.Search{
					{Field: "name", Value: "foo", Operator: searchconditionenum.SearchConditions.CONTAINS},
					{Field: "description", Value: "bar", Operator: searchconditionenum.SearchConditions.CONTAINS},
				},
			},
			res:     "hidden_at IS NULL AND (name ILIKE :search_d0_name_0 OR description ILIKE :search_d0_description_1)",
			resArgs: map[string]any{"search_d0_name_0": "%foo%", "search_d0_description_1": "%bar%"},
		},
		{
			name:  "nested groups",
			where: "hidden_at IS NULL",
			searchGroup: &query.SearchGroup{
				Operator: logicaloperatorenum.LogicalOperators.AND,
				Conditions: []query.Search{
					{Field: "name", Value: "foo", Operator: searchconditionenum.SearchConditions.EQUALS},
				},
				Groups: []*query.SearchGroup{
					{
						Operator: logicaloperatorenum.LogicalOperators.OR,
						Conditions: []query.Search{
							{Field: "slug", Value: "a", Operator: searchconditionenum.SearchConditions.STARTS},
							{Field: "slug", Value: "b", Operator: searchconditionenum.SearchConditions.ENDS},
						},
					},
				},
			},
			res:     "hidden_at IS NULL AND (name = :search_d0_name_0 AND (slug ILIKE :search_d1_slug_0 OR slug ILIKE :search_d1_slug_1))",
			resArgs: map[string]any{"search_d0_name_0": "foo", "search_d1_slug_0": "a%", "search_d1_slug_1": "%b"},
		},

		//validation
		// {
		// 	name:  "unsupported operator errors",
		// 	where: "",
		// 	searchGroup: &query.SearchGroup{
		// 		Operator: logicaloperatorenum.LogicalOperators.AND,
		// 		Conditions: []query.Search{
		// 			{Field: "name", Value: "foo", Operator: searchconditionenum.SearchConditions(99)}, // invalid
		// 		},
		// 	},
		// 	res:   "unsupported search operator",
		// 	isErr: true,
		// },
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			res, args, err := where.BuildNamedWhereSearchesClause(e.where, e.searchGroup)

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
