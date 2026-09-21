package endpoints

import (
	"context"

	"frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/dtos"
	"frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/mappers"
	"frisboo-bank/openapi-generator-service/internal/entities/features/delete_entity/queries"
	entityv1 "frisboo-bank/openapi-generator-service/internal/shared/grpc/gen/entity/v1"
	applicationerror "frisboo-bank/openapi-generator-service/pkg/application_error"
	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	mediatorcontracts "frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

type DeleteEntityEndpoint struct {
	mediator mediatorcontracts.Mediator
}

func NewDeleteEntityEndpoint(mediator mediatorcontracts.Mediator) *DeleteEntityEndpoint {
	validation.AssertNotNil("mediator", mediator)

	return &DeleteEntityEndpoint{
		mediator: mediator,
	}
}

func (e *DeleteEntityEndpoint) DeleteEntity(
	ctx context.Context,
	request *entityv1.DeleteEntityRequest,
) (*entityv1.DeleteEntityResponse, applicationerrorcontracts.AppError) {
	validation.AssertNotNil("request", request)

	query, err := queries.NewDeleteEntityQuery(request)
	if err != nil {
		return nil, applicationerror.NewValidationFailedErrorWrap(ctx, err, nil)
	}

	responseDto, err := mediator.Send[queries.DeleteEntityQuery, dtos.DeleteEntityResponseDto](
		e.mediator,
		ctx,
		query,
	)
	if err != nil {
		if err, ok := applicationerror.IsAppError(err); ok {
			return nil, err
		}
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to delete entity", nil)
	}

	response, err := mappers.MapDeleteEntityResponseToProto(responseDto)
	if err != nil {
		return nil, applicationerror.NewInternalErrorWrap(ctx, err, "failed to map response", nil)
	}

	return response, nil
}
