package models

import (
	"fmt"
	"strings"

	searchconditionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
)

type QuerySearch struct {
	Value    string                              `json:"value"`
	Field    string                              `json:"field"`
	Operator searchconditionEnum.SearchCondition `json:"operator"`
}

func NewQuerySearch(field string, value string, operator searchconditionEnum.SearchCondition) (*QuerySearch, error) {
	query := &QuerySearch{
		Field:    field,
		Value:    value,
		Operator: operator,
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	return query, nil
}

func ParseQuerySearch(field string, value string, operator string) (*QuerySearch, error) {
	op, err := searchconditionEnum.ParseSearchCondition(operator)
	if err != nil {
		return nil, fmt.Errorf("invalid condition %q: %w", op, err)
	}

	return NewQuerySearch(field, value, op)
}

func (qs *QuerySearch) Validate() error {
	if strings.TrimSpace(qs.Value) == "" {
		return fmt.Errorf("search value cannot be empty")
	}

	if strings.TrimSpace(qs.Field) == "" {
		return fmt.Errorf("search field cannot be empty")
	}

	if !qs.Operator.IsValid() {
		return fmt.Errorf("invalid search operator: %s", qs.Operator.String())
	}

	return nil
}
