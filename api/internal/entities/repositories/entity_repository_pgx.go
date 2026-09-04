package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"maps"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	"frisboo-bank/openapi-generator-service/internal/entities/repositories"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/repositories"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/order"
	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/utils/where"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"
	queryModels "frisboo-bank/openapi-generator-service/pkg/query/models"
	filtercomparisonEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/filter_comparison"
	orderdirectionEnum "frisboo-bank/openapi-generator-service/pkg/query/models/enums/order_direction"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/google/uuid"
)

const entityTableName = "entities"

var allowedOrderBy = map[string]struct{}{
	"id":          {},
	"slug":        {},
	"name":        {},
	"description": {},
	"created_at":  {},
	"updated_at":  {},
	"hidden_at":   {},
}

var allowedFilters = map[string]struct{}{
	"hidden_at": {},
}

var _ contracts.EntityRepository = (*EntityRepositoryPgx)(nil)

type EntityRepositoryPgx struct {
	sqlClient sqlclientContracts.SQLXClient
	logger    loggerContracts.Logger
}

func NewEntityRepositoryPgx(
	sqlClient sqlclientContracts.SQLXClient,
	logger loggerContracts.Logger,
) contracts.EntitySQLRepository {
	repo := repositories.NewSQLXRepository( sqlClient, logger, entityTableName)

	return &EntityRepositoryPgx{
		repo: repo,
		logger:    logger,
	}
}

func (e *EntityRepositoryPgx) BeginTx(ctx context.Context) (sqlclientContracts.SQLXTransaction, error) {
	return e.sqlClient.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
}

func (e *EntityRepositoryPgx) CreateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	validation.AssertNotNil("entity", entity)

	query := fmt.Sprintf(`
	  INSERT INTO %s (slug, name, description)
		VALUES (:slug, :name, :description)
		RETURNING id, slug, name, description, version_lock, hidden_at, created_at, updated_at
	`, entityTableName)

	if err := tx.NamedGet(ctx, entity, query, entity); err != nil {
		return nil, fmt.Errorf("insert entity: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	return entity, nil
}

func (e *EntityRepositoryPgx) DeleteEntityByIDTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE id = :id
	`, entityTableName)

	result, err := tx.NamedExec(ctx, query, map[string]any{
		"id": entityID,
	})
	if err != nil {
		return int64(0), fmt.Errorf("delete entity: %w", err)
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		return int64(0), fmt.Errorf("delete entity: %w", err)
	}

	return rowAffected, nil
}

func (e *EntityRepositoryPgx) UpdateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	validation.AssertNotNil("entity", entity)
	validation.AssertNotNil("entity.ID", entity.EntityID)

	currentVersion := entity.VersionLock
	entity.VersionLock++

	query := fmt.Sprintf(`
	  UPDATE %s
	  SET
	    slug = :slug,
	    name = :name,
			description = :description,
			version_lock = :version_lock
    WHERE id = :id
	  AND version_lock = :current_version
	  AND hidden_at IS NULL
		RETURNING id, slug, name, description, version_lock, hidden_at, created_at, updated_at
	`, entityTableName)

	if err := tx.NamedGet(ctx, entity, query, map[string]any{
		"id":              entity.EntityID,
		"slug":            entity.Slug,
		"name":            entity.Name,
		"description":     entity.Description,
		"version_lock":    entity.VersionLock,
		"current_version": currentVersion,
	}); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("entity %s not found or optimistic lock failure", entity.EntityID)
		default:
			return nil, fmt.Errorf("update entity: %w", err)
		}
	}

	return entity, nil
}

func (e *EntityRepositoryPgx) GetEntityById(ctx context.Context, entityID uuid.UUID, filters *query.QueryFilters) (*models.Entity, error) {
	idFilter, err := queryModels.NewQueryFilterWithValue(
		"id",
		[]string{entityID.String()},
		filtercomparisonEnum.FilterComparisons.EQUAL,
	)
	if err != nil {
		return nil, fmt.Errorf("get entity by id: %w", err)
	}
	filters.AddFilter(idFilter)

	whereClause := "WHERE 1=1"

	whereClause, whereArgs, err := where.BuildNamedWhereFiltersClause(whereClause, filters)
	if err != nil {
		return nil, fmt.Errorf("get entity by id: %w", err)
	}

	query := fmt.Sprintf(`
    SELECT id, slug, name, description, version_lock, hidden_at, created_at, updated_at
	  FROM %s
		%s
	`, entityTableName, whereClause)

	var entity models.Entity

	if err := e.sqlClient.NamedGet(ctx, &entity, query, whereArgs); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("entity %s not found", entityID)
		default:
			return nil, fmt.Errorf("get entity by id: %w", err)
		}
	}

	return &entity, nil
}

func (e *EntityRepositoryPgx) GetEntityBySlug(ctx context.Context, slug string, filters *query.QueryFilters) (*models.Entity, error) {
	validation.AssertNotEmpty("slug", slug)

	slugFilter, err := queryModels.NewQueryFilterWithValue(
		"slug",
		[]string{slug},
		filtercomparisonEnum.FilterComparisons.EQUAL,
	)
	if err != nil {
		return nil, fmt.Errorf("get entity by slug: %w", err)
	}
	filters.AddFilter(slugFilter)

	whereClause := "WHERE 1=1"

	whereClause, whereArgs, err := where.BuildNamedWhereFiltersClause(whereClause, filters)
	if err != nil {
		return nil, fmt.Errorf("get entity by slug: %w", err)
	}

	query := fmt.Sprintf(`
    SELECT id, slug, name, description, version_lock, hidden_at, created_at, updated_at
	  FROM %s
		%s
	`, entityTableName, whereClause)

	var entity models.Entity

	if err := e.sqlClient.NamedGet(ctx, &entity, query, whereArgs); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, fmt.Errorf("entity with slug %s not found", slug)
		default:
			return nil, fmt.Errorf("get entity by slug: %w", err)
		}
	}

	return &entity, nil
}

func (e *EntityRepositoryPgx) HideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET hidden_at = now()
		WHERE id = :id
		AND hidden_at IS NULL
	`, entityTableName)

	res, err := e.sqlClient.NamedExec(ctx, query, map[string]any{
		"id": entityID.String(),
	})
	if err != nil {
		return int64(0), fmt.Errorf("hide entity: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return int64(0), fmt.Errorf("hide entity: %w", err)
	}

	return rowsAffected, nil
}

func (e *EntityRepositoryPgx) UnhideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET hidden_at = null
		WHERE id = :id
		AND hidden_at IS NOT NULL
	`, entityTableName)

	res, err := e.sqlClient.NamedExec(ctx, query, map[string]any{
		"id": entityID.String(),
	})
	if err != nil {
		return int64(0), fmt.Errorf("unhide entity: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return int64(0), fmt.Errorf("unhide entity: %w", err)
	}

	return rowsAffected, nil
}

func (e *EntityRepositoryPgx) ListEntitiesTx(
	ctx context.Context,
	tx sqlclientContracts.SQLXTransaction,
	listQuery *query.ListQuery,
) (*query.ListResult[*models.Entity], error) {
	whereClause := "WHERE 1=1"
	var whereArgs map[string]interface{}

  result, err := e.

	var err error

	if filters != nil && !filters.IsEmpty() {
		var filterArgs map[string]interface{}

		whereClause, filterArgs, err = where.BuildNamedWhereFiltersClause(whereClause, filters)
		if err != nil {
			return nil, fmt.Errorf("list entities: %w", err)
		}
		maps.Copy(whereArgs, filterArgs)
	}

	if searches != nil && !searches.IsEmpty() {
		var searchArgs map[string]interface{}

		whereClause, searchArgs, err = where.BuildNamedWhereSearchesClause(whereClause, searches)
		if err != nil {
			return nil, fmt.Errorf("list entities: %w", err)
		}
		maps.Copy(whereArgs, searchArgs)
	}

	if orders == nil || orders.IsEmpty() {
		createdAtOrder, err := queryModels.NewQueryOrder(
			"created_at",
			orderdirectionEnum.OrderDirections.DESC,
		)
		if err != nil {
			return nil, fmt.Errorf("list entities: %w", err)
		}

		orders.AddOrder(createdAtOrder)
	}
	orderClause, err := order.BuildOrderByClause(orders)
	if err != nil {
		return nil, fmt.Errorf("list entities: %w", err)
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
    FROM %s
    %s
	`, entityTableName, whereClause)

	var totalItems int64
	if err := tx.NamedSelect(ctx, &totalItems, countQuery, whereArgs); err != nil {
		return nil, fmt.Errorf("count entities: %w", err)
	}

	var items []*models.Entity

	if totalItems == 0 {
		return query.NewListResult(items, pagination, 0), nil
	}

	whereArgs["limit"] = pagination.Size
	whereArgs["offset"] = pagination.Offset()

	queryString := fmt.Sprintf(`
		SELECT id, slug, name, description, hidden_at, created_at, updated_at
		FROM %s
		%s
		%s
		LIMIT :limit OFFSET :offset
	`, entityTableName, whereClause, orderClause)

	if err := tx.NamedSelect(ctx, &items, queryString, whereArgs); err != nil {
		return nil, fmt.Errorf("select entities: %w", err)
	}

	return query.NewListResult(items, pagination, totalItems, filters, orders, searches), nil
}
