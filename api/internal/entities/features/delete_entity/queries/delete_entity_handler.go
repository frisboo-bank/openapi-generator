package queries

import (
	"context"
	"database/sql"
	"errors"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/dtos"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ mediatorcontracts.RequestHandler[*DeleteEntityQuery, *dtos.DeleteEntityResponseDto] = (*DeleteEntityHandler)(nil)

type DeleteEntityHandler struct {
	repo   contracts.EntityRepository
	logger loggerContracts.Logger
}

func NewDeleteEntityHandler(
	repository contracts.EntityRepository,
	logger loggerContracts.Logger,
) *DeleteEntityHandler {
	return &DeleteEntityHandler{
		repo:   repository,
		logger: logger,
	}
}

func (d *DeleteEntityHandler) Handle(ctx context.Context, request *DeleteEntityQuery) (response *dtos.DeleteEntityResponseDto, err applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	tx, txErr := d.repo.BeginTx(ctx)
	if txErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, txErr, "create transaction failed", nil)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	existingEntity, fetchErr := d.repo.GetEntityBySlugTx(ctx, tx, request.Slug, nil)
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, applicationerror.NewInternalErrorWrap(ctx, fetchErr, "retrieve entity failed", nil)
	}
	if existingEntity == nil {
		return nil, applicationerror.NewNotFoundError(ctx, "entity", request.Slug, nil)
	}

	// deletedEntities, err := d.repo.DeleteEntityByIDTx(ctx, tx, existingEntity.EntityID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to delete entity with error: %w", err)
	// }
	//
	// if deletedEntities != 1 {
	// 	return nil, fmt.Errorf("should have deleted one entity, tried to delete %d", deletedEntities)
	// }
	//
	// if err = tx.Commit(); err != nil {
	// 	return nil, err
	// }

	return &dtos.DeleteEntityResponseDto{}, nil
}
