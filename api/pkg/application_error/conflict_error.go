package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ ConflictError = (*conflictError)(nil)

type ConflictError interface {
	contracts.AppError
	isConflictError()
}

type conflictError struct {
	contracts.AppError
}

func NewConflictError(
	ctx context.Context,
	resourceType string,
	identifier string,
	extraDetails map[string]string,
) ConflictError {
	return NewConflictErrorWrap(ctx, nil, resourceType, identifier, extraDetails)
}

func NewConflictErrorWrap(
	ctx context.Context,
	err error,
	resourceType string,
	identifier string,
	extraDetails map[string]string,
) ConflictError {
	details := map[string]string{
		"resourceType": resourceType,
		"identifier":   identifier,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &conflictError{
		AppError: NewAppError(ctx, contracts.KindConflict, "conflict detected", err, details),
	}
}

func (c *conflictError) isConflictError() {}

func IsConflictError(err error) (ConflictError, bool) {
	if conflictErr, ok := errors.AsType[ConflictError](err); ok {
		return conflictErr, true
	}
	return nil, false
}
