package contracts

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/models"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"

	"github.com/google/uuid"
)

type (
	EntityRepository interface {
		BeginTx(ctx context.Context) (sqlclientContracts.SQLXTransaction, error)
		CreateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
		DeleteEntityByIDTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID) (int64, error)
		GetEntityByID(ctx context.Context, entityID uuid.UUID, query *query.Query) (*models.Entity, error)
		GetEntityBySlug(ctx context.Context, entitySlug string, query *query.Query) (*models.Entity, error)
		GetEntityBySlugTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entitySlug string, query *query.Query) (*models.Entity, error)
		HideEntity(ctx context.Context, entityID uuid.UUID) (int64, error)
		ListEntities(ctx context.Context, query *query.Query) ([]*models.Entity, *query.Pagination, error)
		ListEntitiesTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, query *query.Query) ([]*models.Entity, *query.Pagination, error)
		UnhideEntity(ctx context.Context, entityID uuid.UUID) (int64, error)
		UpdateEntityTx(ctx context.Context, tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
	}

	EntitySQLRepository interface {
		EntityRepository
	}

	EntityCacheRepository interface {
		EntityRepository
	}
)
