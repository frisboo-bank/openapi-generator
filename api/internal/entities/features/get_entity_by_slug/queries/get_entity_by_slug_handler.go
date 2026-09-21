package queries

import (
	"context"
	"database/sql"
	"errors"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/dtos"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ mediatorcontracts.RequestHandler[*GetEntityBySlugQuery, *dtos.GetEntityBySlugResponseDto] = (*GetEntityBySlugHandler)(nil)

type GetEntityBySlugHandler struct {
	repository contracts.EntityRepository
	logger     loggerContracts.Logger
}

func NewGetEntityBySlugHandler(
	repository contracts.EntitySQLRepository,
	logger loggerContracts.Logger,
) *GetEntityBySlugHandler {
	return &GetEntityBySlugHandler{
		repository: repository,
		logger:     logger,
	}
}

func (h *GetEntityBySlugHandler) Handle(ctx context.Context, request *GetEntityBySlugQuery) (*dtos.GetEntityBySlugResponseDto, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	entity, err := h.repository.GetEntityBySlug(ctx, request.Slug, nil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "retrieve entity failed", nil)
	}
	if entity == nil {
		return nil, applicationerror.NewNotFoundError(ctx, "entity", request.Slug, nil)
	}

	return &dtos.GetEntityBySlugResponseDto{
		Entity: entity,
	}, nil
}
