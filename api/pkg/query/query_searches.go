package query

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query/models"

	logicaloperatorEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
)

type QuerySearches struct {
	operator   logicaloperatorEnum.LogicalOperator
	conditions []*models.QuerySearch
	groups     []*QuerySearches
}

func NewQuerySearches(operator logicaloperatorEnum.LogicalOperator) *QuerySearches {
	return &QuerySearches{
		operator:   operator,
		conditions: []*models.QuerySearch{},
		groups:     []*QuerySearches{},
	}
}

func (qs *QuerySearches) Validate(allowedFields map[string]struct{}) error {
	if qs.IsEmpty() {
		return nil
	}

	for _, condition := range qs.conditions {
		if _, ok := allowedFields[condition.Field]; !ok {
			return fmt.Errorf("search field %q not allowed", condition.Field)
		}
	}

	for _, subGroup := range qs.groups {
		if err := subGroup.Validate(allowedFields); err != nil {
			return err
		}
	}

	return nil
}

func (qs *QuerySearches) AddCondition(cond *models.QuerySearch) {
	qs.conditions = append(qs.conditions, cond)
}

func (qs *QuerySearches) AddGroup(group *QuerySearches) {
	qs.groups = append(qs.groups, group)
}

func (qs *QuerySearches) IsEmpty() bool {
	return len(qs.conditions) == 0 && len(qs.groups) == 0
}

func (qs *QuerySearches) GetConditions() []*models.QuerySearch {
	return qs.conditions
}

func (qs *QuerySearches) GetGroups() []*QuerySearches {
	return qs.groups
}

func (qs *QuerySearches) GetOperator() logicaloperatorEnum.LogicalOperator {
	return qs.operator
}
