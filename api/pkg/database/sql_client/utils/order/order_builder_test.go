package order_test

import (
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/order"
	"frisboo-bank/openapi-generator-service/pkg/query"
	orderdirection "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildOrderByClause(t *testing.T) {
	examples := []struct {
		name   string
		order  []query.Order
		result string
		isErr  bool
	}{
		{
			name:   "empty order",
			order:  []query.Order{},
			result: "",
		},
		{
			name:   "order by ASC",
			order:  []query.Order{{Field: "name", Direction: orderdirection.OrderDirections.ASC}},
			result: "ORDER BY name ASC",
		},
		{
			name:   "order by DESC",
			order:  []query.Order{{Field: "created_at", Direction: orderdirection.OrderDirections.DESC}},
			result: "ORDER BY created_at DESC",
		},
		{
			name: "multiple order by",
			order: []query.Order{
				{Field: "name", Direction: orderdirection.OrderDirections.ASC},
				{Field: "created_at", Direction: orderdirection.OrderDirections.DESC},
				{Field: "updated_at", Direction: orderdirection.OrderDirections.ASC},
			},
			result: "ORDER BY name ASC, created_at DESC, updated_at ASC",
		},
		{
			name:   "order by field with prefix",
			order:  []query.Order{{Field: "customer.name", Direction: orderdirection.OrderDirections.ASC}},
			result: "ORDER BY customer.name ASC",
		},
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			res, err := order.BuildOrderByClause(e.order)

			if e.isErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, e.result, res)
		})
	}
}
