package applicationerror

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
	"frisboo-bank/openapi-generator-service/pkg/constants"
)

var _ contracts.AppError = (*appError)(nil)

type appError struct {
	kind          contracts.ErrorKind
	message       string
	requestID     string
	correlationID string
	details       map[string]string
	cause         error
	timestamp     time.Time
}

func NewAppError(
	ctx context.Context,
	kind contracts.ErrorKind,
	message string,
	cause error,
	extraDetails map[string]string,
) contracts.AppError {
	details := make(map[string]string)
	if extraDetails != nil {
		maps.Copy(details, extraDetails)
	}

	requestID := ""
	if rID, ok := ctx.Value(constants.HeaderRequestIDKey).(string); ok {
		requestID = rID
	}

	correlationID := ""
	if cID, ok := ctx.Value(constants.CorrelationIDKey).(string); ok {
		correlationID = cID
	}

	return &appError{
		kind:          kind,
		message:       message,
		requestID:     requestID,
		correlationID: correlationID,
		details:       details,
		cause:         cause,
		timestamp:     time.Now().UTC(),
	}
}

func (e *appError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.kind, e.message, e.cause)
	}
	return fmt.Sprintf("[%s] %s", e.kind, e.message)
}

func (e *appError) Kind() contracts.ErrorKind  { return e.kind }
func (e *appError) Message() string            { return e.message }
func (e *appError) CorrelationID() string      { return e.correlationID }
func (e *appError) RequestID() string          { return e.requestID }
func (e *appError) Details() map[string]string { return e.details }
func (e *appError) Unwrap() error              { return e.cause }
func (e *appError) Timestamp() time.Time       { return e.timestamp }

func IsAppError(err error) (contracts.AppError, bool) {
	if appError, ok := errors.AsType[contracts.AppError](err); ok {
		return appError, true
	}
	return nil, false
}
