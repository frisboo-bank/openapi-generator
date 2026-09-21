package endpoints

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/mappers"
	"frisboo-bank/openapi-generator-service/internal/entities/features/update_entity/queries"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type UpdateEntityEndpoint struct {
	mediator mediatorcontracts.Mediator
}

func NewUpdateEntityEndpoint(mediator mediatorcontracts.Mediator) *UpdateEntityEndpoint {
	validation.AssertNotNil("mediator", mediator)

	return &UpdateEntityEndpoint{
		mediator: mediator,
	}
}

func (e *UpdateEntityEndpoint) UpdateEntity(
	ctx context.Context,
	request *entityv1.UpdateEntityRequest,
) (*entityv1.UpdateEntityResponse, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	query, err := queries.NewUpdateEntityQuery(request)
	if err != nil {
		return nil, applicationerror.NewValidationFailedErrorWrap(ctx, err, nil)
	}

	responseDto, err := mediator.Send[queries.UpdateEntityQuery, dtos.UpdateEntityResponseDto](e.mediator, ctx, query)
	if err != nil {
		if err, ok := applicationerror.IsAppError(err); ok {
			return nil, err
		}
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to update entity", nil)
	}

	return mappers.MapUpdateEntityResponseToProto(responseDto), nil
}
