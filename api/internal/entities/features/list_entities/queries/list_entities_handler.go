package queries

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	listdtos "frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/dtos"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

var _ mediatorcontracts.RequestHandler[*ListEntitiesQuery, *listdtos.ListEntitiesResponseDto] = (*ListEntitiesHandler)(nil)

type ListEntitiesHandler struct {
	repo   contracts.EntityRepository
	logger loggercontracts.Logger
}

func NewListEntitiesHandler(
	repository contracts.EntityRepository,
	logger loggercontracts.Logger,
) *ListEntitiesHandler {
	validation.AssertNotNil("repository", repository)
	validation.AssertNotNil("logger", logger)

	return &ListEntitiesHandler{
		repo:   repository,
		logger: logger,
	}
}

func (h *ListEntitiesHandler) Handle(ctx context.Context, request *ListEntitiesQuery) (*listdtos.ListEntitiesResponseDto, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	entities, pagination, err := h.repo.ListEntities(ctx, &request.Query)
	if err != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "list entities failed", nil)
	}

	return &listdtos.ListEntitiesResponseDto{
		Entities:   entities,
		Pagination: pagination,
	}, nil
}
