package order

import (
	"fmt"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/query"
	orderdirection "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
)

func BuildOrderByClause(orders []query.Order) (string, error) {
	if len(orders) == 0 {
		return "", nil
	}

	parts := make([]string, len(orders))

	for i, o := range orders {
		direction := "ASC"
		if o.Direction == orderdirection.OrderDirections.DESC {
			direction = "DESC"
		}
		parts[i] = fmt.Sprintf("%s %s", o.Field, direction)
	}

	return "ORDER BY " + strings.Join(parts, ", "), nil
}
