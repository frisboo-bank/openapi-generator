package mappers

import (
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/dtos"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
)

func MapListEntitiesResponseToProto(responseDto *dtos.ListEntitiesResponseDto) (*entityv1.ListEntitiesResponse, error) {
	return &entityv1.ListEntitiesResponse{}, nil
	// if len(responseDto.Entities) == 0 {
	// 	return &entityv1.ListEntitiesResponse{
	// 		Data:       []*entityv1.Entity{},
	// 		Pagination: &commonv1.PaginationResponse{},
	// 	}, nil
	// }
	//
	// data := make([]*entityv1.Entity, 0, len(responseDto.Entities))
	// for _, entity := range responseDto.Entities {
	// 	data = append(data, &entityv1.Entity{
	// 		Slug:        entity.Slug,
	// 		Name:        entity.Name,
	// 		Description: utils.StringToWrapper(&entity.Description),
	// 		VersionLock: entity.VersionLock,
	// 		HiddenAt:    utils.TimeToTimestamp(entity.HiddenAt),
	// 		CreatedAt:   utils.TimeToTimestamp(entity.CreatedAt),
	// 		UpdatedAt:   utils.TimeToTimestamp(entity.UpdatedAt),
	// 	})
	// }
	//
	// return &entityv1.ListEntitiesResponse{
	// 	Data: data,
	// 	Pagination: &commonv1.PaginationResponse{
	// 		Page:       responseDto.Pagination.Page,
	// 		PageSize:   responseDto.Pagination.PageSize,
	// 		TotalCount: responseDto.Pagination.TotalCount,
	// 		PageCount:  responseDto.Pagination.PageCount,
	// 	},
	// }, nil
}
