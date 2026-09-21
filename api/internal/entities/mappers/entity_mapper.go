package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/utils"
)

func MapEntityToProto(entity *models.Entity) (*entityv1.Entity, error) {
	return &entityv1.Entity{
		Slug:        entity.Slug,
		Name:        entity.Name,
		VersionLock: entity.VersionLock,
		HiddenAt:    utils.TimeToTimestamp(entity.HiddenAt),
		CreatedAt:   utils.TimeToTimestamp(entity.CreatedAt),
		UpdatedAt:   utils.TimeToTimestamp(entity.UpdatedAt),
	}, nil
}
