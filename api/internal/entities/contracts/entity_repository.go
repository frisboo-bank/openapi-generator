package contracts

import (
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	sqlclientContracts "frisboo-bank/openapi-generator-service/pkg/database/sql_client/contracts"
	"frisboo-bank/openapi-generator-service/pkg/query"

	"github.com/google/uuid"
)

type (
	EntityRepository interface {
		BeginTx() (sqlclientContracts.SQLXTransaction, error)
		CreateEntityTx(tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
		DeleteEntityByIDTx(tx sqlclientContracts.SQLXTransaction, entityID uuid.UUID) (int64, error)
		GetEntityByID(entityID uuid.UUID, query *query.Query) (*models.Entity, error)
		GetEntityBySlug(entitySlug string, query *query.Query) (*models.Entity, error)
		GetEntityBySlugTx(tx sqlclientContracts.SQLXTransaction, entitySlug string, query *query.Query) (*models.Entity, error)
		HideEntity(entityID uuid.UUID) (int64, error)
		ListEntities(query *query.Query) ([]*models.Entity, *query.Pagination, error)
		ListEntitiesTx(tx sqlclientContracts.SQLXTransaction, query *query.Query) ([]*models.Entity, *query.Pagination, error)
		UnhideEntity(entityID uuid.UUID) (int64, error)
		UpdateEntityTx(tx sqlclientContracts.SQLXTransaction, entity *models.Entity) (*models.Entity, error)
	}

	EntitySQLRepository interface {
		EntityRepository
	}

	EntityCacheRepository interface {
		EntityRepository
	}
)
