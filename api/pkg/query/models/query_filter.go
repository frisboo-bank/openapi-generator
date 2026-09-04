package models

import (
	"fmt"
	"strings"

	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
)

type QueryFilter struct {
	Field      string                                `json:"field"`
	Values     []string                              `json:"values"`
	Comparator filtercomparisonEnum.FilterComparison `json:"comparison"`
}

func newQueryFilter(field string, values []string, comparator filtercomparisonEnum.FilterComparison) (*QueryFilter, error) {
	query := &QueryFilter{
		Field:      field,
		Values:     values,
		Comparator: comparator,
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	return query, nil
}

func NewQueryFilter(field string, comp filtercomparisonEnum.FilterComparison) (*QueryFilter, error) {
	if comp.RequiresValue() {
		return nil, fmt.Errorf("comparison %s requires a value, use NewFilterWithValueModel", comp.String())
	}
	return newQueryFilter(field, nil, comp)
}

func NewQueryFilterWithValue(field string, values []string, comp filtercomparisonEnum.FilterComparison) (*QueryFilter, error) {
	if !comp.RequiresValue() {
		return nil, fmt.Errorf("comparison %s does not accept a value, use NewQueryFilter", comp.String())
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("at least one value required for comparison %s", comp.String())
	}
	return newQueryFilter(field, values, comp)
}

func (f *QueryFilter) Validate() error {
	if strings.TrimSpace(f.Field) == "" {
		return fmt.Errorf("filter field must not be empty")
	}

	if !f.Comparator.IsValid() {
		return fmt.Errorf("invalid filter comparison")
	}

	return f.Comparator.ValidateValue(f.Values)
}
