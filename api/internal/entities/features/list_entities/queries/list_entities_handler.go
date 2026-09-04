package queries

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/contracts"
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/dtos"
	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
)

type ListEntitiesHandler struct {
	repo   contracts.EntitySQLRepository
	logger loggerContracts.Logger
}

func NewListEntitiesHandler(
	repo contracts.EntitySQLRepository,
	logger loggerContracts.Logger,
) *ListEntitiesHandler {
	return &ListEntitiesHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *ListEntitiesHandler) Handle(
	ctx context.Context,
	query *ListEntries,
) (*dtos.ListEntitiesResponseDto, error) {
	tx, err := h.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.repo.ListEntitiesTx(
		ctx, tx,
		query.ListQuery,
	)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return &dtos.ListEntitiesResponseDto{
		Items: result,
		// Pagination: result
		// Filters:    result
		// Orders:     result
		// Searches:   result.S
	}, nil
}
