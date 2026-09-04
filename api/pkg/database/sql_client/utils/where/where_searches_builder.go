package where

import (
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/query"
)

func BuildNamedWhereSearchesClause(
	where string,
	searches *query.QuerySearches,
) (string, map[string]interface{}, error) {
	var conditions []string
	var args map[string]interface{}

	if searches.IsEmpty() {
		return where, args, nil
	}

	if len(conditions) == 0 {
		return where, args, nil
	}

	return where + " AND " + strings.Join(conditions, " AND "), args, nil
}
