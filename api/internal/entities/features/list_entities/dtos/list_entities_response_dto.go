package dtos

import (
	"frisboo-bank/openapi-generator-service/internal/entities/dtos"
	"frisboo-bank/openapi-generator-service/pkg/query"
)

type ListEntitiesResponseDto struct {
	Items      []*dtos.EntityDto
	Pagination *query.Pagination
	TotalItems int64
}
