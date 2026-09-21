package applicationerror

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

var _ RateLimitedError = (*rateLimitedError)(nil)

type RateLimitedError interface {
	contracts.AppError
	isRateLimitedError()
}

type rateLimitedError struct {
	contracts.AppError
}

func NewRateLimitedError(
	ctx context.Context,
	limit string,
	window string,
	retryAfterSeconds int64,
	extraDetails map[string]string,
) RateLimitedError {
	return NewRateLimitedErrorWrap(ctx, nil, limit, window, retryAfterSeconds, extraDetails)
}

func NewRateLimitedErrorWrap(
	ctx context.Context,
	err error,
	limit string,
	window string,
	retryAfterSeconds int64,
	extraDetails map[string]string,
) RateLimitedError {
	details := map[string]string{
		"limit":             limit,
		"window":            window,
		"retryAfterSeconds": fmt.Sprintf("%d", retryAfterSeconds),
	}
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	return &rateLimitedError{
		AppError: NewAppError(ctx, contracts.KindRateLimited, "rate limited", err, details),
	}
}

func (r *rateLimitedError) isRateLimitedError() {}

func IsRateLimitedError(err error) (RateLimitedError, bool) {
	if rateLimitedError, ok := errors.AsType[RateLimitedError](err); ok {
		return rateLimitedError, true
	}
	return nil, false
}
