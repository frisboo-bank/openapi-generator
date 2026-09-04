package query

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query/models"
)

type QueryOrders struct {
	root []*models.QueryOrder
}

func NewQueryOrders(orderBy []*models.QueryOrder) *QueryOrders {
	return &QueryOrders{
		root: orderBy,
	}
}

func (qf *QueryOrders) Validate(allowedFields map[string]struct{}) error {
	if qf.IsEmpty() {
		return nil
	}

	for i, o := range qf.root {
		if o == nil {
			return fmt.Errorf("order model at index %d is nil", i)
		}
		if err := o.Validate(); err != nil {
			return fmt.Errorf("order model at index %d: %w", i, err)
		}
		if _, ok := allowedFields[o.Field]; !ok {
			return fmt.Errorf("order field %q is not allowed", o.Field)
		}
	}

	return nil
}

func (qf *QueryOrders) AddOrder(o *models.QueryOrder) {
	qf.root = append(qf.root, o)
}

func (qo *QueryOrders) IsEmpty() bool {
	return len(qo.root) == 0
}

func (qo *QueryOrders) GetOrders() []*models.QueryOrder {
	return qo.root
}
