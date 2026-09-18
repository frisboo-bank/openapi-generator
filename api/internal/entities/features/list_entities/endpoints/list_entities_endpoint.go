package endpoints

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/mappers"
	"frisboo-bank/openapi-generator-service/internal/entities/features/list_entities/queries"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type ListEntitiesEndpoint struct {
	mediator mediatorcontracts.Mediator
}

func NewListEntitiesEndpoint(mediator mediatorcontracts.Mediator) *ListEntitiesEndpoint {
	validation.AssertNotNil("mediator", mediator)

	return &ListEntitiesEndpoint{
		mediator: mediator,
	}
}

func (e *ListEntitiesEndpoint) ListEntities(
	ctx context.Context,
	request *entityv1.ListEntitiesRequest,
) (*entityv1.ListEntitiesResponse, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	query, err := queries.NewListEntitiesQuery(request)
	if err != nil {
		return nil, applicationerror.NewValidationFailedErrorWrap(ctx, err, nil)
	}

	responseDto, err := mediator.Send[queries.ListEntitiesQuery, dtos.ListEntitiesResponseDto](
		e.mediator,
		ctx,
		query,
	)
	if err != nil {
		if err, ok := applicationerror.IsAppError(err); ok {
			return nil, err
		}
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to fetch entity list", nil)
	}

	response, err := mappers.MapListEntitiesResponseToProto(responseDto)
	if err != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to map response", nil)
	}

	return response, nil
}
