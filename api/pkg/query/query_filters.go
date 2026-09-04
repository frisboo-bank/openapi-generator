package query

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query/models"
)

type QueryFilters struct {
	root []*models.QueryFilter
}

func NewQueryFilters(filters []*models.QueryFilter) *QueryFilters {
	return &QueryFilters{
		root: filters,
	}
}

func (qf *QueryFilters) Validate(allowedFields map[string]struct{}) error {
	if qf.IsEmpty() {
		return nil
	}

	for i, f := range qf.root {
		if f == nil {
			return fmt.Errorf("filter model at index %d is nil", i)
		}

		if err := f.Validate(); err != nil {
			return fmt.Errorf("filter model at index %d: %w", i, err)
		}

		if _, ok := allowedFields[f.Field]; !ok {
			return fmt.Errorf("filter field %q is not allowed", f.Field)
		}
	}

	return nil
}

func (qf *QueryFilters) AddFilter(f *models.QueryFilter) {
	qf.root = append(qf.root, f)
}

func (qf *QueryFilters) IsEmpty() bool {
	return len(qf.root) == 0
}

func (qf *QueryFilters) GetFilters() []*models.QueryFilter {
	return qf.root
}
