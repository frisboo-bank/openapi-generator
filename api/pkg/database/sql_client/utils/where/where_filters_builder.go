package where

import (
	"fmt"
	"maps"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/mappers"
	"frisboo-bank/openapi-generator-service/pkg/query"
	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
)

func BuildNamedWhereFiltersClause(where string, filters []query.Filter) (string, map[string]any, error) {
	if len(filters) == 0 {
		return where, nil, nil
	}

	var conditions []string
	args := make(map[string]any)

	for i, f := range filters {
		var condition string
		var conditionArgs map[string]any
		var err error

		switch f.Comparator {
		case filtercomparisonEnum.FilterComparisons.BETWEEN:
			condition, conditionArgs, err = betweenCondition(i, f)
		case filtercomparisonEnum.FilterComparisons.IN:
			condition, conditionArgs, err = inCondition(i, f)
		default:
			if f.Comparator.Options().RequiresValue {
				condition, conditionArgs, err = unaryCondition(i, f)
			} else {
				condition, conditionArgs, err = nullableCondition(f)
			}
		}
		if err != nil {
			return "", nil, fmt.Errorf("filter %q: %w", f.Field, err)
		}

		conditions = append(conditions, condition)
		maps.Copy(args, conditionArgs)
	}

	if len(conditions) == 0 {
		return where, args, nil
	}

	clause := strings.Join(conditions, " AND ")
	if where == "" {
		return clause, args, nil
	}
	return where + " AND " + clause, args, nil
}

func betweenCondition(idx int, f query.Filter) (string, map[string]any, error) {
	if len(f.Values) != 2 {
		return "", nil, fmt.Errorf("BETWEEN requires exactly 2 values")
	}

	lowKey := fmt.Sprintf("%s_%d_low", f.Field, idx)
	highKey := fmt.Sprintf("%s_%d_high", f.Field, idx)

	return fmt.Sprintf("%s BETWEEN :%s AND :%s", f.Field, lowKey, highKey), map[string]any{
		lowKey:  f.Values[0],
		highKey: f.Values[1],
	}, nil
}

func inCondition(idx int, f query.Filter) (string, map[string]any, error) {
	if len(f.Values) == 0 {
		return "", nil, fmt.Errorf("IN requires at least 1 value")
	}

	keys := make([]string, len(f.Values))
	args := make(map[string]any, len(f.Values))

	for i, v := range f.Values {
		key := fmt.Sprintf("%s_%d_%d", f.Field, idx, i)
		keys[i] = ":" + key
		args[key] = v
	}

	return fmt.Sprintf("%s IN (%s)", f.Field, strings.Join(keys, ", ")), args, nil
}

func unaryCondition(idx int, f query.Filter) (string, map[string]any, error) {
	if len(f.Values) != 1 {
		return "", nil, fmt.Errorf("%s requires exactly 1 value", f.Comparator.String())
	}

	key := fmt.Sprintf("%s_%d", f.Field, idx)
	operation, err := mappers.FilterComparisonToSQLOperator(f.Comparator)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("%s %s :%s", f.Field, operation, key),
		map[string]any{key: f.Values[0]},
		nil
}

func nullableCondition(f query.Filter) (string, map[string]any, error) {
	operation, err := mappers.FilterComparisonToSQLOperator(f.Comparator)
	if err != nil {
		return "", nil, err
	}

	return fmt.Sprintf("%s %s", f.Field, operation), nil, nil
}
