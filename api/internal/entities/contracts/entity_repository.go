package contracts

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/models"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"

	"github.com/google/uuid"
)

type EntityRepository interface {
	BeginTx(ctx context.Context) (sqlclientContracts.SQLXTransaction, error)
	CreateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
	DeleteEntityByIDTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID) (int64, error)
	UpdateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
	HideEntity(ctx context.Context, entityID uuid.UUID) (int64, error)
	UnhideEntity(ctx context.Context, entityID uuid.UUID) (int64, error)
	ListEntitiesTx(
		ctx context.Context,
		tx sqlclientContracts.SQLXTransaction,
		listQuery *query.ListQuery,
	) (*query.ListResult[*models.Entity], error)
	GetEntityById(
		ctx context.Context,
		entityID uuid.UUID,
		filters *query.QueryFilters,
	) (*models.Entity, error)
	GetEntityBySlug(
		ctx context.Context,
		entitySlug string,
		filters *query.QueryFilters,
	) (*models.Entity, error)
}

type EntitySQLRepository interface {
	EntityRepository
}
