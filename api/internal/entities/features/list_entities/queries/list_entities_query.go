package queries

import (
	"fmt"
	"frisboo-bank/openapi-generator-service/internal/entities/models"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	"frisboo-bank/openapi-generator-service/pkg/core/mappers"
	"frisboo-bank/openapi-generator-service/pkg/query"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type ListEntitiesQuery struct {
	Query query.Query
}

func NewListEntitiesQuery(request *entityv1.ListEntitiesRequest) (*ListEntitiesQuery, error) {
	validation.AssertNotNil("request", request)
	validation.AssertNotNil("request.Body", request.Body)

	q, err := mappers.MapProtoV1ToQuery(
		request.Body.GetFilters(),
		request.Body.GetOrders(),
		request.Body.GetSearches(),
		request.Body.GetPagination(),
	)
	if err != nil {
		return nil, fmt.Errorf("map request to query: %w", err)
	}

	cmd := &ListEntitiesQuery{Query: q}
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	return cmd, nil
}

func (l *ListEntitiesQuery) Validate() error {
	return l.Query.Validate(query.AllowedQueryFieldsFor[models.Entity]())
}
