package dtos

import (
	"frisboo-bank/openapi-generator-service/pkg/query"
)

type ListEntitiesRequestDto struct {
	pagination *query.Pagination
	filters    *query.QueryFilters
	orders     *query.QueryOrders
	searches   *query.QuerySearches
}
