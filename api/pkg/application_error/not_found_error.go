package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ NotFoundError = (*notFoundError)(nil)

type NotFoundError interface {
	contracts.AppError
	isNotFoundError()
}

type notFoundError struct {
	contracts.AppError
}

func NewNotFoundError(
	ctx context.Context,
	resourceType string,
	identifier string,
	extraDetails map[string]string,
) NotFoundError {
	return NewNotFoundErrorWrap(ctx, nil, resourceType, identifier, extraDetails)
}

func NewNotFoundErrorWrap(
	ctx context.Context,
	err error,
	resourceType string,
	identifier string,
	extraDetails map[string]string,
) NotFoundError {
	details := map[string]string{
		"resourceType": resourceType,
		"identifier":   identifier,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &notFoundError{
		AppError: NewAppError(ctx, contracts.KindNotFound, "not found", err, details),
	}
}

func (n *notFoundError) isNotFoundError() {}

func IsNotFoundError(err error) (NotFoundError, bool) {
	if notFoundError, ok := errors.AsType[NotFoundError](err); ok {
		return notFoundError, true
	}
	return nil, false
}
