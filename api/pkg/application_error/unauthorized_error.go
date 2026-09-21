package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ UnauthorizedError = (*unauthorizedError)(nil)

type UnauthorizedError interface {
	contracts.AppError
	isUnauthorizedError()
}

type unauthorizedError struct {
	contracts.AppError
}

func NewUnauthorizedError(
	ctx context.Context,
	reason string,
	extraDetails map[string]string,
) UnauthorizedError {
	return NewUnauthorizedErrorWrap(ctx, nil, reason, extraDetails)
}

func NewUnauthorizedErrorWrap(
	ctx context.Context,
	err error,
	reason string,
	extraDetails map[string]string,
) UnauthorizedError {
	details := map[string]string{
		"reason": reason,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &unauthorizedError{
		AppError: NewAppError(ctx, contracts.KindUnauthorized, reason, err, details),
	}
}

func (u *unauthorizedError) isUnauthorizedError() {}

func IsUnauthorizedError(err error) (UnauthorizedError, bool) {
	if unauthorizedError, ok := errors.AsType[UnauthorizedError](err); ok {
		return unauthorizedError, true
	}
	return nil, false
}
