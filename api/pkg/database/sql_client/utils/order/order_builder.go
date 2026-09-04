package order

import (
	"fmt"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/query"
)

func BuildOrderByClause(orders *query.QueryOrders) (string, error) {
	if orders.IsEmpty() {
		return "", nil
	}

	orderModels := orders.GetOrders()
	parts := make([]string, 0, len(orderModels))

	for i, m := range orderModels {
		if m == nil {
			return "", fmt.Errorf("order model at index %d is nil", i)
		}

		parts = append(parts, m.Field+" "+m.Direction.String())
	}

	if len(parts) == 0 {
		return "", nil
	}

	return "ORDER BY " + strings.Join(parts, ", "), nil
}
