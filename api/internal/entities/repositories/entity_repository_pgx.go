package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/google/uuid"
)

const entityTableName = "entities"

var maxPaginationLimit = 100

var _ contracts.EntitySQLRepository = (*EntityRepositoryPgx)(nil)

type EntityRepositoryPgx struct {
	sqlClient sqlclientContracts.SQLXClientAdapter
	logger    loggerContracts.Logger
}

func NewEntityRepositoryPgx(
	sqlClient sqlclientContracts.SQLXClientAdapter,
	logger loggerContracts.Logger,
) contracts.EntitySQLRepository {
	validation.AssertNotNil("sqlClient", sqlClient)
	validation.AssertNotNil("logger", logger)

	return &EntityRepositoryPgx{
		sqlClient: sqlClient,
		logger:    logger,
	}
}

func (e *EntityRepositoryPgx) BeginTx(ctx context.Context) (sqlclientContracts.SQLXTransaction, error) {
	return e.sqlClient.BeginTransaction(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

func (e *EntityRepositoryPgx) CreateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	validation.AssertNotNil("entity", entity)

	query := fmt.Sprintf(`
    INSERT INTO %s (slug, name, description)
    VALUES (:slug, :name, :description)
    RETURNING id, slug, name, description, version_lock, hidden_at, created_at, updated_at
  `, entityTableName)

	if err := tx.NamedGet(ctx, entity, query, map[string]any{
		"slug":        entity.Slug,
		"name":        entity.Name,
		"description": entity.Description,
	}); err != nil {
		return nil, fmt.Errorf("insert entity: %w", err)
	}
	return entity, nil
}

func (e *EntityRepositoryPgx) DeleteEntityByIDTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = :id`, entityTableName)

	result, err := tx.NamedExec(ctx, query, map[string]any{"id": entityID})
	if err != nil {
		return 0, fmt.Errorf("delete entity: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *EntityRepositoryPgx) GetEntityByID(ctx context.Context, entityID uuid.UUID, query *query.Query) (*models.Entity, error) {
	return e.doGetEntityByID(ctx, nil, entityID, query)
}

func (e *EntityRepositoryPgx) GetEntityBySlug(ctx context.Context, entitySlug string, query *query.Query) (*models.Entity, error) {
	return e.doGetEntityBySlug(ctx, nil, entitySlug, query)
}

func (e *EntityRepositoryPgx) GetEntityBySlugTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entitySlug string, query *query.Query) (*models.Entity, error) {
	return e.doGetEntityBySlug(ctx, tx, entitySlug, query)
}

func (e *EntityRepositoryPgx) HideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET hidden_at = now()
		WHERE id = :id
		AND hidden_at IS NULL
	`, entityTableName)

	res, err := e.sqlClient.NamedExec(ctx, query, map[string]any{"id": entityID})
	if err != nil {
		return 0, fmt.Errorf("hide entity: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *EntityRepositoryPgx) ListEntities(ctx context.Context, query *query.Query) ([]*models.Entity, *query.Pagination, error) {
	return e.doListEntities(ctx, nil, query)
}

func (e *EntityRepositoryPgx) ListEntitiesTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, query *query.Query) ([]*models.Entity, *query.Pagination, error) {
	return e.doListEntities(ctx, tx, query)
}

func (e *EntityRepositoryPgx) UnhideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	query := fmt.Sprintf(`
		UPDATE %s
		SET hidden_at = null
		WHERE id = :id
		AND hidden_at IS NOT NULL
	`, entityTableName)

	res, err := e.sqlClient.NamedExec(ctx, query, map[string]any{"id": entityID})
	if err != nil {
		return 0, fmt.Errorf("unhide entity: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *EntityRepositoryPgx) UpdateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	validation.AssertNotNil("entity", entity)
	validation.AssertNotNil("entity.EntityID", entity.EntityID)

	currentVersion := entity.VersionLock
	entity.VersionLock++

	query := fmt.Sprintf(`
		UPDATE %s
		SET name = :name, description = :description, version_lock = :version_lock
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("entity not found or optimistic lock failure")
		}
		return nil, fmt.Errorf("update entity: %w", err)
	}

	return entity, nil
}

func (e *EntityRepositoryPgx) doGetEntityByID(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID, q *query.Query) (*models.Entity, error) {
	validation.AssertNotEmpty("entityID", entityID.String())
	validation.AssertNotNil("q", q)

	return nil, nil
}

func (e *EntityRepositoryPgx) doGetEntityBySlug(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entitySlug string, q *query.Query) (*models.Entity, error) {
	panic("unimplemented")
}

func (e *EntityRepositoryPgx) doListEntities(ctx context.Context, tx sqlclientContracts.SQLXTransaction, q *query.Query) ([]*models.Entity, *query.Pagination, error) {
	panic("unimplemented")
}
