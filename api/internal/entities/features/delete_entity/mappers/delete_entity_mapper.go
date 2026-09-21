package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/dtos"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
)

func MapDeleteEntityResponseToProto(dto *dtos.DeleteEntityResponseDto) (*entityv1.DeleteEntityResponse, error) {
	return &entityv1.DeleteEntityResponse{}, nil
}
