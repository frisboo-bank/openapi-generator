package dtos

import (
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	"frisboo-bank/openapi-generator-service/pkg/query"
)

type ListEntitiesResponseDto struct {
	Entities   []*models.Entity
	Pagination *query.Pagination
}
