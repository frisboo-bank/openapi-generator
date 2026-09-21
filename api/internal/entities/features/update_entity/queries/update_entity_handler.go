package queries

import (
	"context"
	"database/sql"
	"errors"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/dtos"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ mediatorcontracts.RequestHandler[*UpdateEntityQuery, *dtos.UpdateEntityResponseDto] = (*UpdateEntityHandler)(nil)

type UpdateEntityHandler struct {
	repo   contracts.EntityRepository
	logger loggerContracts.Logger
}

func NewUpdateEntityHandler(
	repository contracts.EntityRepository,
	logger loggerContracts.Logger,
) *UpdateEntityHandler {
	validation.AssertNotNil("repository", repository)
	validation.AssertNotNil("logger", logger)

	return &UpdateEntityHandler{
		repo:   repository,
		logger: logger,
	}
}

// func (h *ListEntitiesHandler) Handle(ctx context.Context, request *ListEntitiesQuery) (*listdtos.ListEntitiesResponseDto, applicationerrorcontracts.AppError) {
// 	validation.AssertNotNil("request", request)
//
// 	entities, pagination, err := h.repo.ListEntities(ctx, &request.Query)
// 	if err != nil {
// 		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "list entities failed", nil)
// 	}
//
// 	return &listdtos.ListEntitiesResponseDto{
// 		Entities:   entities,
// 		Pagination: pagination,
// 	}, nil
// }

func (h *UpdateEntityHandler) Handle(ctx context.Context, request *UpdateEntityQuery) (response *dtos.UpdateEntityResponseDto, err applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	tx, txErr := h.repo.BeginTx(ctx)
	if txErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, txErr, "create transaction failed", nil)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	existingEntity, fetchErr := h.repo.GetEntityBySlugTx(ctx, tx, request.Slug, nil)
	if fetchErr != nil && !errors.Is(fetchErr, sql.ErrNoRows) {
		return nil, applicationerror.NewInternalErrorWrap(ctx, fetchErr, "retrieve entity failed", nil)
	}
	if existingEntity == nil {
		return nil, applicationerror.NewNotFoundError(ctx, "entity", request.Slug, nil)
	}

	if request.VersionLock != existingEntity.VersionLock {
		return nil, applicationerror.NewVersionMismatchError(
			ctx,
			"entity",
			request.Slug,
			existingEntity.VersionLock,
			request.VersionLock,
			nil,
		)
	}

	if request.Name != nil {
		existingEntity.Name = *request.Name
	}
	if request.Description != nil {
		existingEntity.Description = *request.Description
	}

	updatedEntity, updateErr := h.repo.UpdateEntityTx(ctx, tx, existingEntity)
	if updateErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, updateErr, "update entity failed", nil)
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, commitErr, "commit transaction failed", nil)
	}

	return &dtos.UpdateEntityResponseDto{
		Slug:        updatedEntity.Slug,
		Name:        updatedEntity.Name,
		Description: updatedEntity.Description,
		VersionLock: updatedEntity.VersionLock,
		HiddenAt:    updatedEntity.HiddenAt,
		CreatedAt:   updatedEntity.CreatedAt,
		UpdatedAt:   updatedEntity.UpdatedAt,
	}, nil
}
