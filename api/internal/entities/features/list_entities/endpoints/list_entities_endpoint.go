package endpoints

import (
	"context"

	entityv1 "frisboo-bank/openapi-generator-service/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/queries"
)

var _ entityv1.EntityServiceServer = (*ListEntitiesEndpoint)(nil)

type ListEntitiesEndpoint struct {
	entityv1.UnimplementedEntityServiceServer
	handler *queries.ListEntitiesHandler
}

func NewListEntitiesEndpoint(handler *queries.ListEntitiesHandler) *ListEntitiesEndpoint {
	return &ListEntitiesEndpoint{
		handler: handler,
	}
}

func (e *ListEntitiesEndpoint) ListEntities(
	ctx context.Context,
	req *entityv1.ListEntitiesRequest,
) (*entityv1.ListEntitiesResponse, error) {
	return nil, nil
}
