package where

import (
	"fmt"
	"frisboo-bank/openapi-generator-service/pkg/query"
	logicaloperatorenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
	searchconditionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
	"maps"
	"strings"
)

func BuildNamedWhereSearchesClause(where string, searchGroup *query.SearchGroup) (string, map[string]any, error) {
	if searchGroup == nil || searchGroup.IsEmpty() {
		return where, nil, nil
	}

	clause, args, err := buildSearchGroup(0, searchGroup)
	if err != nil {
		return "", nil, err
	}

	if where == "" {
		return clause, args, nil
	}
	return where + " AND (" + clause + ")", args, nil
}

func buildSearchGroup(depth int, searchGroup *query.SearchGroup) (string, map[string]any, error) {
	if searchGroup == nil || searchGroup.IsEmpty() {
		return "", nil, nil
	}

	var conditions []string
	args := make(map[string]any)

	for i, c := range searchGroup.Conditions {
		key := fmt.Sprintf("search_d%d_%s_%d", depth, c.Field, i)

		var condition string
		var value string

		switch c.Operator {
		case searchconditionenum.SearchConditions.EQUALS:
			condition = fmt.Sprintf("%s = :%s", c.Field, key)
			value = c.Value
		case searchconditionenum.SearchConditions.CONTAINS:
			condition = fmt.Sprintf("%s ILIKE :%s", c.Field, key)
			value = "%" + c.Value + "%"
		case searchconditionenum.SearchConditions.STARTS:
			condition = fmt.Sprintf("%s ILIKE :%s", c.Field, key)
			value = c.Value + "%"
		case searchconditionenum.SearchConditions.ENDS:
			condition = fmt.Sprintf("%s ILIKE :%s", c.Field, key)
			value = "%" + c.Value
		default:
			return "", nil, fmt.Errorf("unsupported search operator: %s", c.Operator.String())
		}

		conditions = append(conditions, condition)
		args[key] = value
	}

	for i, g := range searchGroup.Groups {
		if g == nil {
			continue
		}

		subCondition, subArgs, err := buildSearchGroup(depth+1, g)
		if err != nil {
			return "", nil, fmt.Errorf("subgroup[%d]: %w", i, err)
		}
		if subCondition == "" {
			continue
		}
		conditions = append(conditions, "("+subCondition+")")
		maps.Copy(args, subArgs)
	}

	if len(conditions) == 0 {
		return "", nil, nil
	}

	op := "AND"
	if searchGroup.Operator == logicaloperatorenum.LogicalOperators.OR {
		op = "OR"
	}

	return strings.Join(conditions, " "+op+" "), args, nil
}
