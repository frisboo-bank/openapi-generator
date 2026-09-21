package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ InternalError = (*internalError)(nil)

type InternalError interface {
	contracts.AppError
	isInternalError()
}

type internalError struct {
	contracts.AppError
}

func NewInternalError(
	ctx context.Context,
	message string,
	extraDetails map[string]string,
) InternalError {
	return NewInternalErrorWrap(ctx, nil, message, extraDetails)
}

func NewInternalErrorFromService(
	ctx context.Context,
	cause error,
	service string,
	operation string,
	extraDetails map[string]string,
) InternalError {
	details := map[string]string{
		"service":   service,
		"operation": operation,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return NewInternalErrorWrap(ctx, cause, "dependency failed", details)
}

func NewInternalErrorWrap(
	ctx context.Context,
	err error,
	message string,
	extraDetails map[string]string,
) InternalError {
	details := make(map[string]string)
	if err != nil {
		details["cause"] = err.Error()
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &internalError{
		AppError: NewAppError(ctx, contracts.KindInternal, message, err, details),
	}
}

func (i *internalError) isInternalError() {}

func IsInternalError(err error) (InternalError, bool) {
	if internalError, ok := errors.AsType[InternalError](err); ok {
		return internalError, true
	}
	return nil, false
}
