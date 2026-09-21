package queries

import (
	"context"
	"database/sql"
	"errors"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ mediatorcontracts.RequestHandler[*CreateEntityQuery, *dtos.CreateEntityResponseDto] = (*CreateEntityHandler)(nil)

type CreateEntityHandler struct {
	repo   contracts.EntityRepository
	logger loggercontracts.Logger
}

func NewCreateEntityHandler(repository contracts.EntityRepository, logger loggercontracts.Logger) *CreateEntityHandler {
	return &CreateEntityHandler{
		repo:   repository,
		logger: logger,
	}
}

func (c *CreateEntityHandler) Handle(ctx context.Context, request *CreateEntityQuery) (response *dtos.CreateEntityResponseDto, err applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	tx, txErr := c.repo.BeginTx(ctx)
	if txErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, txErr, "create transaction failed", nil)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	existingEntity, fetchErr := c.repo.GetEntityBySlugTx(ctx, tx, request.Slug, nil)
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, applicationerror.NewInternalErrorWrap(ctx, fetchErr, "retrieve entity failed", nil)
	}
	if existingEntity != nil {
		return nil, applicationerror.NewConflictError(ctx, "entity", request.Slug, nil)
	}

	entity := &models.Entity{
		Slug:        request.Slug,
		Name:        request.Name,
		Description: *request.Description,
	}

	createdEntity, createErr := c.repo.CreateEntityTx(ctx, tx, entity)
	if createErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, createErr, "create entity failed", nil)
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, commitErr, "commit transaction failed", nil)
	}

	return &dtos.CreateEntityResponseDto{
		Slug:        createdEntity.Slug,
		Name:        createdEntity.Name,
		Description: &createdEntity.Description,
		VersionLock: createdEntity.VersionLock,
		HiddenAt:    createdEntity.HiddenAt,
		CreatedAt:   createdEntity.CreatedAt,
		UpdatedAt:   createdEntity.UpdatedAt,
	}, nil
}
