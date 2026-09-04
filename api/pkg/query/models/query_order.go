package models

import (
	"fmt"

	orderdirectionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	"frisboo-bank/openapi-generator-service/pkg/utils"
)

type QueryOrder struct {
	Field     string                            `json:"field"`
	Direction orderdirectionEnum.OrderDirection `json:"direction"`
}

func NewQueryOrder(field string, direction orderdirectionEnum.OrderDirection) (*QueryOrder, error) {
	query := &QueryOrder{
		Field:     field,
		Direction: direction,
	}

	if err := query.Validate(); err != nil {
		return nil, err
	}

	return query, nil
}

func ParseQueryOrder(field string, direction string) (*QueryOrder, error) {
	dir, err := orderdirectionEnum.ParseOrderDirection(direction)
	if err != nil {
		return nil, fmt.Errorf("invalid direction %q: %w", direction, err)
	}

	return NewQueryOrder(field, dir)
}

func (om *QueryOrder) Validate() error {
	if utils.StripSpace(om.Field) == "" {
		return fmt.Errorf("order field must not be empty")
	}

	if !om.Direction.IsValid() {
		return fmt.Errorf("order %q direction is not set", om.Field)
	}

	return nil
}
