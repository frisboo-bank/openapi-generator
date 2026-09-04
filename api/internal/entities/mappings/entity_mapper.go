package mappings

import (
	"frisboo-bank/openapi-generator-service/internal/entities/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
)

func MapEntityToEntityDto(entity *models.Entity) *dtos.EntityDto {
	return &dtos.EntityDto{
		EntityID:    entity.EntityID,
		Slug:        entity.Slug,
		Name:        entity.Name,
		Description: entity.Description,
		VersionLock: entity.VersionLock,
		HiddenAt:    entity.HiddenAt,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
