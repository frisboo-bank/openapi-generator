package query

import (
	"encoding/json"
	"fmt"
)

type ListResult[T any] struct {
	Items      []T         `json:"items"`
	Pagination *Pagination `json:"pagination"`
	TotalItems int64       `json:"total_items"`
	TotalPages int64       `json:"total_pages"`
}

func NewListResult[T any](
	items []T,
	pagination *Pagination,
	totalItems int64,
) *ListResult[T] {
	if items == nil {
		items = []T{}
	}
	return &ListResult[T]{
		Items:      items,
		Pagination: pagination,
		TotalItems: totalItems,
		TotalPages: pagination.TotalPages(totalItems),
	}
}

func (p *ListResult[T]) String() string {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf("<error marshaling ListResult: %v>", err)
	}
	return string(data)
}
