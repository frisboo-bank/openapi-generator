package endpoints

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/mappers"
	"frisboo-bank/openapi-generator-service/internal/entities/features/get_entity_by_slug/queries"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type GetEntityBySlugEndpoint struct {
	mediator mediatorcontracts.Mediator
}

func NewGetEntityBySlugEndpoint(mediator mediatorcontracts.Mediator) *GetEntityBySlugEndpoint {
	validation.AssertNotNil("mediator", mediator)

	return &GetEntityBySlugEndpoint{
		mediator: mediator,
	}
}

func (e *GetEntityBySlugEndpoint) GetEntityBySlug(
	ctx context.Context,
	request *entityv1.GetEntityBySlugRequest,
) (*entityv1.GetEntityBySlugResponse, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	query, err := queries.NewGetEntityBySlugQuery(request)
	if err != nil {
		return nil, applicationerror.NewValidationFailedErrorWrap(ctx, err, nil)
	}

	responseDto, err := mediator.Send[queries.GetEntityBySlugQuery, dtos.GetEntityBySlugResponseDto](
		e.mediator,
		ctx,
		query,
	)
	if err != nil {
		if err, ok := applicationerror.IsAppError(err); ok {
			return nil, err
		}
		return nil, applicationerror.NewInternalErrorWrap(ctx, err,
			fmt.Sprintf("failed to fetch entity by slug `%s`", request.Slug), nil)
	}

	response, err := mappers.MapGetEntityBySlugResponseToProto(responseDto)
	if err != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to map response", nil)
	}

	return response, nil
}
