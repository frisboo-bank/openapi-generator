package applicationerror

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ VersionMismatchError = (*versionMismatchError)(nil)

type VersionMismatchError interface {
	contracts.AppError
	isVersionMismatchError()
}

type versionMismatchError struct {
	contracts.AppError
}

func NewVersionMismatchError(
	ctx context.Context,
	resourceType string,
	identifier string,
	expectedVersion int64,
	actualVersion int64,
	extraDetails map[string]string,
) VersionMismatchError {
	return NewVersionMismatchErrorWrap(ctx, nil, resourceType, identifier, expectedVersion, actualVersion, extraDetails)
}

func NewVersionMismatchErrorWrap(
	ctx context.Context,
	err error,
	resourceType string,
	identifier string,
	expectedVersion int64,
	actualVersion int64,
	extraDetails map[string]string,
) VersionMismatchError {
	details := map[string]string{
		"resourceType":    resourceType,
		"identifier":      identifier,
		"expectedVersion": fmt.Sprintf("%d", expectedVersion),
		"actualVersion":   fmt.Sprintf("%d", actualVersion),
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &versionMismatchError{
		AppError: NewAppError(ctx, contracts.KindVersionMismatch, "version mismatch", err, details),
	}
}

func (v *versionMismatchError) isVersionMismatchError() {}

func IsVersionMismatchError(err error) (VersionMismatchError, bool) {
	if versionMismatchErr, ok := errors.AsType[VersionMismatchError](err); ok {
		return versionMismatchErr, true
	}
	return nil, false
}
