package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	vendorValidation "github.com/go-ozzo/ozzo-validation"
)

var _ ValidationFailedError = (*validationFailedError)(nil)

type ValidationFailedError interface {
	contracts.AppError
	isValidationFailedError()
}

type validationFailedError struct {
	contracts.AppError
}

func NewValidationFailedError(
	ctx context.Context,
	field string,
	violation string,
	extraDetails map[string]string,
) ValidationFailedError {
	details := map[string]string{
		"field":     field,
		"violation": violation,
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &validationFailedError{
		AppError: NewAppError(ctx, contracts.KindValidationFailed, "validation failed", nil, details),
	}
}

func NewValidationFailedErrorWrap(
	ctx context.Context,
	err error,
	extraDetails map[string]string,
) ValidationFailedError {
	details := make(map[string]string)

	if vErr, ok := err.(vendorValidation.Errors); ok {
		maps.Copy(details, validation.FlattenErrors(vErr))
	}

	if len(details) == 0 {
		details["error"] = err.Error()
	}

	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &validationFailedError{
		AppError: NewAppError(ctx, contracts.KindValidationFailed, "validation failed", err, details),
	}
}

func (v *validationFailedError) isValidationFailedError() {}

func IsValidationFailedError(err error) (ValidationFailedError, bool) {
	if validationError, ok := errors.AsType[ValidationFailedError](err); ok {
		return validationError, true
	}
	return nil, false
}
