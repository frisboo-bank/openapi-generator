package applicationerror

import (
	"context"
	"errors"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ BusinessRuleError = (*businessRuleError)(nil)

type BusinessRuleError interface {
	contracts.AppError
	isBusinessRuleError()
}

type businessRuleError struct {
	contracts.AppError
}

func NewBusinessRuleError(
	ctx context.Context,
	rule string,
	message string,
	extraDetails map[string]string,
) BusinessRuleError {
	return NewBusinessRuleErrorWrap(ctx, nil, rule, message, extraDetails)
}

func NewBusinessRuleErrorWrap(
	ctx context.Context,
	err error,
	rule string,
	message string,
	extraDetails map[string]string,
) BusinessRuleError {
	details := map[string]string{
		"rule": rule,
	}

	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &businessRuleError{
		AppError: NewAppError(ctx, contracts.KindBusinessRule, message, err, details),
	}
}

func (b *businessRuleError) isBusinessRuleError() {}

func IsBusinessRuleError(err error) (BusinessRuleError, bool) {
	if businessError, ok := errors.AsType[BusinessRuleError](err); ok {
		return businessError, true
	}
	return nil, false
}
