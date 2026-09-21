package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/dtos"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
)

func MapCreateEntityResponseToProto(dto *dtos.CreateEntityResponseDto) (*entityv1.CreateEntityResponse, error) {
	return &entityv1.CreateEntityResponse{
		// Entity: &entityv1.Entity{
		// 	Slug:        dto.Slug,
		// 	Name:        dto.Name,
		// 	Description: utils.StringToWrapper(dto.Description),
		// 	VersionLock: dto.VersionLock,
		// 	HiddenAt:    utils.TimeToTimestamp(dto.HiddenAt),
		// 	CreatedAt:   utils.TimeToTimestamp(dto.CreatedAt),
		// 	UpdatedAt:   utils.TimeToTimestamp(dto.UpdatedAt),
		// },
	}, nil
}
