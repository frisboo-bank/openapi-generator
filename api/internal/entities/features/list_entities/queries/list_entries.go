package queries

import "frisboo-bank/openapi-generator-service/pkg/query"

type ListEntries struct {
	*query.ListQuery
}

func NewListEntries(query *query.ListQuery) *ListEntries {
	return &ListEntries{
		ListQuery: query,
	}
}
