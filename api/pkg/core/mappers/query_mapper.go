package mappers

import (
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/query"
	filtercomparisonenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
	logicaloperatorenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/logical_operator"
	orderdirectionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	searchconditionenum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/search_condition"
	queryv1 "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/resources/grpc/gen/query/v1"
)

func MapProtoV1ToQuery(
	filters *queryv1.QueryFilters,
	orders *queryv1.QueryOrders,
	searches *queryv1.QuerySearches,
	pagination *queryv1.PaginationRequest,
) (query.Query, error) {
	q := query.Query{}

	if filters != nil {
		f, err := mapProtoV1ToFilters(filters)
		if err != nil {
			return query.Query{}, fmt.Errorf("filters: %w", err)
		}
		q.Filters = f
	}

	if orders != nil {
		o, err := mapProtoV1ToOrders(orders)
		if err != nil {
			return query.Query{}, fmt.Errorf("orders: %w", err)
		}
		q.Orders = o
	}

	if searches != nil {
		s, err := mapProtoV1ToSearches(searches)
		if err != nil {
			return query.Query{}, fmt.Errorf("searches: %w", err)
		}
		q.Searches = s
	}

	if pagination != nil {
		q.Pagination = mapProtoV1ToPagination(pagination)
	}

	return q, nil
}

func mapProtoV1ToFilters(filters *queryv1.QueryFilters) ([]query.Filter, error) {
	res := make([]query.Filter, 0, len(filters.GetFilters()))

	for i, f := range filters.GetFilters() {
		cmp, err := filtercomparisonenum.ParseFilterComparison(f.GetComparison())
		if err != nil {
			return nil, fmt.Errorf("[%d] field %q: invalid comparison %q: %w", i, f.GetField(), f.GetComparison(), err)
		}

		res = append(res, query.Filter{
			Field:      f.GetField(),
			Values:     f.GetValues(),
			Comparator: cmp,
		})
	}

	return res, nil
}

func mapProtoV1ToOrders(orders *queryv1.QueryOrders) ([]query.Order, error) {
	res := make([]query.Order, 0, len(orders.GetOrders()))

	for i, o := range orders.GetOrders() {
		dir, err := orderdirectionenum.ParseOrderDirection(o.GetDirection())
		if err != nil {
			return nil, fmt.Errorf("[%d] field %q: invalid direction %q: %w", i, o.GetField(), o.GetDirection(), err)
		}

		res = append(res, query.Order{
			Field:     o.GetField(),
			Direction: dir,
		})
	}

	return res, nil
}

func mapProtoV1ToSearches(searches *queryv1.QuerySearches) (*query.SearchGroup, error) {
	if searches == nil {
		return nil, nil
	}

	if len(searches.GetConditions()) == 0 &&
		len(searches.GetGroups()) == 0 &&
		searches.GetOperator() == queryv1.LogicalOperator_LOGICAL_OPERATOR_UNSPECIFIED {
		return nil, nil
	}

	protoOp := searches.GetOperator()
	op, err := logicaloperatorenum.ParseLogicalOperator(int32(protoOp))
	if err != nil {
		return nil, fmt.Errorf("invalid operator %v: %w", protoOp, err)
	}

	sg := &query.SearchGroup{
		Operator: op,
	}

	for i, c := range searches.GetConditions() {
		protoCond := c.GetCondition()
		cOp, err := searchconditionenum.ParseSearchCondition(int32(protoCond))
		if err != nil {
			return nil, fmt.Errorf("condition[%d] field %q: invalid condition %v: %w", i, c.GetField(), protoCond, err)
		}

		sg.Conditions = append(sg.Conditions, query.Search{
			Field:    c.GetField(),
			Value:    c.GetValue(),
			Operator: cOp,
		})
	}

	for i, g := range searches.GetGroups() {
		child, err := mapProtoV1ToSearches(g)
		if err != nil {
			return nil, fmt.Errorf("group[%d]: %w", i, err)
		}
		sg.Groups = append(sg.Groups, child)
	}

	return sg, nil
}

func mapProtoV1ToPagination(pagination *queryv1.PaginationRequest) *query.Pagination {
	return &query.Pagination{
		Cursor: pagination.GetCursor(),
		Limit:  pagination.GetLimit(),
	}
}
