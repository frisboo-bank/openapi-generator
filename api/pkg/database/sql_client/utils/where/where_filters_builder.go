package where

import (
	"fmt"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/mappers"
	"frisboo-bank/openapi-generator-service/pkg/query"
	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
)

func BuildNamedWhereFiltersClause(where string, filters *query.QueryFilters) (string, map[string]interface{}, error) {
	var conditions []string
	var args map[string]interface{}

	if filters.IsEmpty() {
		return where, args, nil
	}

	for i, filter := range filters.GetFilters() {
		if filter == nil {
			return "", nil, fmt.Errorf("filter at index %d is nil", i)
		}

		if err := filter.Validate(); err != nil {
			return "", nil, fmt.Errorf("filter at index %d: %w", i, err)
		}

		operator, err := mappers.FilterComparisonToSQLOperator(filter.Comparator)
		if err != nil {
			return "", nil, err
		}

		var condition string
		var conditionArgs map[string]interface{}

		switch filter.Comparator {
		case filtercomparisonEnum.FilterComparisons.BETWEEN:
			condition, conditionArgs, err = betweenCondition(filter)
		case filtercomparisonEnum.FilterComparisons.IN:
			condition, conditionArgs, err = inCondition(filter)
		default:
			if filter.Comparator.RequiresValue() {
				condition, conditionArgs, err = unaryCondition(filter, operator)
			} else {
				condition = fmt.Sprintf("%s %s", filter.Field, operator)
			}
		}
		if err != nil {
			return "", nil, fmt.Errorf("filter %q: %w", filter.Field, err)
		}

		conditions = append(conditions, condition)

		for k, v := range conditionArgs {
			if _, exists := args[k]; exists {
				return "", nil, fmt.Errorf("duplicate placeholder %q from filter %q", k, filter.Field)
			}
			args[k] = v
		}
	}

	if len(conditions) == 0 {
		return where, args, nil
	}

	return where + " AND " + strings.Join(conditions, " AND "), args, nil
}
