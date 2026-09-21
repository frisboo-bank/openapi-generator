package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/dtos"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/utils"
)

func MapGetEntityBySlugResponseToProto(
	dto *dtos.GetEntityBySlugResponseDto,
) (*entityv1.GetEntityBySlugResponse, error) {
	return &entityv1.GetEntityBySlugResponse{
		Data: &entityv1.GetEntityBySlugResponse_Data{
			Entity: &entityv1.Entity{
				Slug:        dto.Entity.Slug,
				Name:        dto.Entity.Name,
				Description: &dto.Entity.Description,
				VersionLock: dto.Entity.VersionLock,
				HiddenAt:    utils.TimeToTimestamp(dto.Entity.HiddenAt),
				CreatedAt:   utils.TimeToTimestamp(dto.Entity.CreatedAt),
				UpdatedAt:   utils.TimeToTimestamp(dto.Entity.UpdatedAt),
			},
		},
	}, nil
}
