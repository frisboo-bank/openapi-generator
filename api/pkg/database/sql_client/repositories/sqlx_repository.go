package repositories

import (
	"context"
	"fmt"
	"maps"
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/order"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/where"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type SQLXRepository[TModel interface{}] struct {
	sqlClient contracts.SQLXClient
	logger    loggerContracts.Logger
	tableName string
}

func NewSQLXRepository[TModel interface{}](
	sqlClient contracts.SQLXClient,
	logger loggerContracts.Logger,
	tableName string,
) *SQLXRepository[TModel] {
	return &SQLXRepository[TModel]{
		sqlClient: sqlClient,
		logger:    logger,
		tableName: tableName,
	}
}

func (s *SQLXRepository[TModel]) PaginatedList(
	ctx context.Context,
	tx contracts.SQLXTransaction,
	pagination *query.Pagination,
	filters *query.QueryFilters,
	orders *query.QueryOrders,
	searches *query.QuerySearches,
	fields []string,
) (*query.ListResult[TModel], error) {
	validation.AssertNotNil("filters", filters)
	validation.AssertNotNil("orders", orders)
	validation.AssertNotNil("searches", searches)
	validation.Assert(len(fields) > 0, "at least one field required")

	var err error

	whereClause := "WHERE 1=1"
	var whereArgs map[string]interface{}

	if !filters.IsEmpty() {
		var filterArgs map[string]interface{}
		whereClause, filterArgs, err = where.BuildNamedWhereFiltersClause(whereClause, filters)
		if err != nil {
			return nil, fmt.Errorf("build filters: %w", err)
		}
		maps.Copy(whereArgs, filterArgs)
	}

	if !searches.IsEmpty() {
		var searchArgs map[string]interface{}
		whereClause, searchArgs, err = where.BuildNamedWhereSearchesClause(whereClause, searches)
		if err != nil {
			return nil, fmt.Errorf("build searches: %w", err)
		}
		maps.Copy(whereArgs, searchArgs)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
    FROM %s
    %s
	`, s.tableName, whereClause)

	var totalItems int64
	if err := tx.NamedSelect(ctx, &totalItems, countQuery, whereArgs); err != nil {
		return nil, fmt.Errorf("count: %w", err)
	}

	if totalItems == 0 {
		return query.NewListResult([]TModel{}, pagination, 0), nil
	}

	orderClause := ""
	if !orders.IsEmpty() {
		orderClause, err = order.BuildOrderByClause(orders)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}
	}

	var paginationClause string
	if pagination != nil {
		paginationClause = "LIMIT :limit OFFSET :offset"
		whereArgs["limit"] = pagination.Size
		whereArgs["offset"] = pagination.Offset()
	}

	queryString := fmt.Sprintf(
		`
		SELECT %s
		FROM %s
		%s
		%s
		%s
	`,
		strings.Join(fields, ", "),
		s.tableName,
		whereClause,
		orderClause,
		paginationClause,
	)

	var items []TModel
	if err := tx.NamedSelect(ctx, &items, queryString, whereArgs); err != nil {
		return nil, fmt.Errorf("select entities: %w", err)
	}

	return query.NewListResult(
		items,
		pagination,
		totalItems,
	), nil
}
