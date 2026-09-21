package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/dtos"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/utils"
)

func MapUpdateEntityResponseToProto(dto *dtos.UpdateEntityResponseDto) *entityv1.UpdateEntityResponse {
	return &entityv1.UpdateEntityResponse{
		Data: &entityv1.UpdateEntityResponse_Data{
			Entity: &entityv1.Entity{
				Slug:        dto.Slug,
				Name:        dto.Name,
				Description: &dto.Description,
				VersionLock: dto.VersionLock,
				HiddenAt:    utils.TimeToTimestamp(dto.HiddenAt),
				CreatedAt:   utils.TimeToTimestamp(dto.CreatedAt),
				UpdatedAt:   utils.TimeToTimestamp(dto.UpdatedAt),
			},
		},
	}
}
