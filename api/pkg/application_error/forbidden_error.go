package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ ForbiddenError = (*forbiddenError)(nil)

type ForbiddenError interface {
	contracts.AppError
	isForbiddenError()
}

type forbiddenError struct {
	contracts.AppError
}

func NewForbiddenError(
	ctx context.Context,
	resource string,
	action string,
	principal string,
	extraDetails map[string]string,
) ForbiddenError {
	return NewForbiddenErrorWrap(ctx, nil, resource, action, principal, extraDetails)
}

func NewForbiddenErrorWrap(
	ctx context.Context,
	err error,
	resource string,
	action string,
	principal string,
	extraDetails map[string]string,
) ForbiddenError {
	details := map[string]string{
		"resource":  resource,
		"action":    action,
		"principal": principal,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &forbiddenError{
		AppError: NewAppError(ctx, contracts.KindForbidden, "forbidden", err, details),
	}
}

func (f *forbiddenError) isForbiddenError() {}

func IsForbiddenError(err error) (ForbiddenError, bool) {
	if forbiddenError, ok := errors.AsType[ForbiddenError](err); ok {
		return forbiddenError, true
	}
	return nil, false
}
