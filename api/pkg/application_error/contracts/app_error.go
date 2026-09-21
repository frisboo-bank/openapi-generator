package contracts

type ErrorKind string

const (
	KindBusinessRule     ErrorKind = "BUSINESS_RULE_VIOLATION"
	KindConflict         ErrorKind = "CONFLICT"
	KindForbidden        ErrorKind = "FORBIDDEN"
	KindInternal         ErrorKind = "INTERNAL"
	KindNotFound         ErrorKind = "NOT_FOUND"
	KindRateLimited      ErrorKind = "RATE_LIMITED"
	KindUnauthorized     ErrorKind = "UNAUTHORIZED"
	KindValidationFailed ErrorKind = "VALIDATION_FAILED"
	KindVersionMismatch  ErrorKind = "VERSION_MISMATCH"
)

type AppError interface {
	error
	Kind() ErrorKind
	Message() string
	RequestID() string
	CorrelationID() string
	Details() map[string]string
	Unwrap() error
}
