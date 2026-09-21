package endpoints

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/mappers"
	"frisboo-bank/openapi-generator-service/internal/entities/features/create_entity/queries"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type CreateEntityEndpoint struct {
	mediator mediatorcontracts.Mediator
}

func NewCreateEntityEndpoint(mediator mediatorcontracts.Mediator) *CreateEntityEndpoint {
	validation.AssertNotNil("mediator", mediator)

	return &CreateEntityEndpoint{
		mediator: mediator,
	}
}

func (e *CreateEntityEndpoint) CreateEntity(
	ctx context.Context,
	request *entityv1.CreateEntityRequest,
) (*entityv1.CreateEntityResponse, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	query, err := queries.NewCreateEntityQuery(request)
	if err != nil {
		return nil, applicationerror.NewValidationFailedErrorWrap(ctx, err, nil)
	}

	responseDto, err := mediator.Send[queries.CreateEntityQuery, dtos.CreateEntityResponseDto](e.mediator, ctx, query)
	if err != nil {
		if err, ok := applicationerror.IsAppError(err); ok {
			return nil, err
		}
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to create entity", nil)
	}

	response, err := mappers.MapCreateEntityResponseToProto(responseDto)
	if err != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to map response", nil)
	}

	return response, nil
}
