package repositories

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	cachecontracts "frisboo-bank/openapi-generator-service/pkg/cache/contracts"
	sqlclientcontracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"

	"github.com/google/uuid"
)

var _ contracts.EntityCacheRepository = (*EntityRepositoryRedis)(nil)

const (
	cacheKeyEntityID   = "entity:id:%s"
	cacheKeyEntitySlug = "entity:slug:%s"
)

type EntityRepositoryRedis struct {
	cache    cachecontracts.Cache
	delegate contracts.EntitySQLRepository
	logger   loggercontracts.Logger
}

func NewEntityRepositoryRedis(
	cache cachecontracts.Cache,
	delegate contracts.EntitySQLRepository,
	logger loggercontracts.Logger,
) contracts.EntityCacheRepository {
	return &EntityRepositoryRedis{
		cache:    cache,
		delegate: delegate,
		logger:   logger,
	}
}

// BeginTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) BeginTx(ctx context.Context) (sqlclientcontracts.SQLXTransaction, error) {
	panic("unimplemented")
}

// CreateEntityTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) CreateEntityTx(ctx context.Context, tx sqlclientcontracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	panic("unimplemented")
}

// DeleteEntityByIDTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) DeleteEntityByIDTx(ctx context.Context, tx sqlclientcontracts.SQLXTransaction, entityID uuid.UUID) (int64, error) {
	panic("unimplemented")
}

// GetEntityByID implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) GetEntityByID(ctx context.Context, entityID uuid.UUID, query *query.Query) (*models.Entity, error) {
	panic("unimplemented")
}

// GetEntityBySlug implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) GetEntityBySlug(ctx context.Context, entitySlug string, query *query.Query) (*models.Entity, error) {
	panic("unimplemented")
}

// GetEntityBySlugTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) GetEntityBySlugTx(ctx context.Context, tx sqlclientcontracts.SQLXTransaction, entitySlug string, query *query.Query) (*models.Entity, error) {
	panic("unimplemented")
}

// HideEntity implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) HideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	panic("unimplemented")
}

// ListEntities implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) ListEntities(ctx context.Context, query *query.Query) ([]*models.Entity, *query.Pagination, error) {
	panic("unimplemented")
}

// ListEntitiesTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) ListEntitiesTx(ctx context.Context, tx sqlclientcontracts.SQLXTransaction, query *query.Query) ([]*models.Entity, *query.Pagination, error) {
	panic("unimplemented")
}

// UnhideEntity implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) UnhideEntity(ctx context.Context, entityID uuid.UUID) (int64, error) {
	panic("unimplemented")
}

// UpdateEntityTx implements [contracts.EntityCacheRepository].
func (e *EntityRepositoryRedis) UpdateEntityTx(ctx context.Context, tx sqlclientcontracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error) {
	panic("unimplemented")
}
