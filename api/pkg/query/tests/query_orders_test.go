package query_test

import (
	"testing"

	"frisboo-bank/openapi-generator-service/pkg/query"
	orderdirection "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"

	"pgregory.net/rapid"
)

func generateQueryOrder() *rapid.Generator[*query.Order] {
	return rapid.Custom(func(t *rapid.T) *query.Order {
		var field string
		for field == "" {
			field = rapid.String().Draw(t, "field")
		}

		directions := make([]orderdirection.OrderDirection, 0)
		for direction := range orderdirection.OrderDirections.All() {
			if direction.IsValid() {
				directions = append(directions, direction)
			}
		}
		direction := rapid.SampledFrom(directions).Draw(t, "direction")

		order, err := query.NewQueryOrder(field, direction)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		return order
	})
}

func TestQueryOrders(t *testing.T) {
	t.Run("empty orders is always valid", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			allowed := rapid.SliceOf(rapid.String()).Draw(t, "allowed")

			qo := query.NewQueryOrders()

			if err := qo.Validate(allowed); err != nil {
				t.Fatalf("empty orders should validate: %v", err)
			}
		})
	})

	t.Run("orders included not allowed field should fail validation", func(t *testing.T) {
		rapid.Check(t, func(t *rapid.T) {
			orders := rapid.SliceOfN(generateQueryOrder(), 2, 20).Draw(t, "orders")

			badIdx := rapid.IntRange(0, len(orders)-1).Draw(t, "badIdx")
			notAllowed := orders[badIdx].Field

			allowedSet := make(map[string]struct{})
			for i, order := range orders {
				if i == badIdx {
					continue
				}
				allowedSet[order.Field] = struct{}{}
			}
			allowed := make([]string, 0, len(allowedSet))
			for field := range allowedSet {
				allowed = append(allowed, field)
			}

			qo := query.NewQueryOrders(orders...)

			err := qo.Validate(allowed)
			if err == nil {
				t.Fatalf("expected validation error for field %q not in allowed %q", notAllowed, allowed)
			}
		})
	})
}
