package query

import (
	filtercomparisonenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
	logicaloperatorenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
	orderdirectionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	searchconditionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
)

const (
	DefaultPaginationLimit = 20
	MaxPaginationLimit     = 100
)

type Query struct {
	Filters    []Filter
	Orders     []Order
	Searches   *SearchGroup
	Pagination *Pagination
}

type Filter struct {
	Field      string                                `json:"field"`
	Values     []string                              `json:"values"`
	Comparator filtercomparisonenum.FilterComparison `json:"comparison"`
}

type Order struct {
	Field     string                            `json:"field"`
	Direction orderdirectionenum.OrderDirection `json:"direction"`
}

type SearchGroup struct {
	Operator   logicaloperatorenum.LogicalOperator `json:"operator"`
	Conditions []Search                            `json:"conditions"`
	Groups     []*SearchGroup                      `json:"groups"`
}

func (sg *SearchGroup) IsEmpty() bool {
	if sg == nil {
		return true
	}
	return len(sg.Conditions) == 0 && len(sg.Groups) == 0
}

type Search struct {
	Field    string                              `json:"field"`
	Value    string                              `json:"value"`
	Operator searchconditionenum.SearchCondition `json:"operator"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
	Limit  int32  `json:"limit"`
}
