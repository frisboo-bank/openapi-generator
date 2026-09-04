package syserrors

import (
	verr "gitlab.com/tozd/go/errors"
)

type E = verr.E

// New creates a new error with a stack trace and optional initial details.
func New(message string, kv ...any) E {
	return verr.WithDetails(verr.New(message), kv...)
}

// Newf creates a formatted error with a stack trace.
func Newf(format string, args ...any) E {
	return verr.Errorf(format, args...)
}

// Wrap wraps an existing error with a new message and attaches optional details.
// Returns nil if err is nil.
func Wrap(err error, message string, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(verr.Wrap(err, message), ensureDetails(kv...)...)
}

// Wrapf is the formatted variant of Wrap.
func Wrapf(err error, format string, args ...any) E {
	return verr.Wrapf(err, format, args...)
}

// Message adds a prefix message and key-value details to an error.
func Message(err error, prefix []string, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(verr.WithMessage(err, prefix...), kv...)
}

// Messagef adds a formatted prefix message and key-value details to an error.
func Messagef(err error, format string, args []any, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(
		verr.WithMessagef(err, format, args...),
		kv...,
	)
}

// WithDetails adds key-value details to an error.
func WithDetails(err error, kv ...any) E {
	return verr.WithDetails(err, kv...)
}

// Join joins errors (skips nil)
func Join(errs ...error) E {
	return verr.Join(errs...)
}

// Prefix adds prefix to an error
func Prefix(err error, prefix ...error) E {
	return verr.Prefix(err, prefix...)
}

// WithStack ensures the error has a stack
func WithStack(err error) E {
	return verr.WithStack(err)
}

// WrapWith wraps an error with another error
func WrapWith(err error, with error) E {
	return verr.WrapWith(err, with)
}

// Cause returns the underlying cause of an error
func Cause(err error) error {
	return verr.Cause(err)
}

// AllDetails returns the aggregated map
func AllDetails(err error) map[string]any {
	return verr.AllDetails(err)
}

func Get(err error, key string) (value any, exists bool) {
	if err == nil {
		return nil, false
	}
	details := verr.AllDetails(err)
	v, ok := details[key]
	return v, ok
}

// ensureDetails validates key/value detail pairs length.
func ensureDetails(kv ...any) []any {
	if len(kv)%2 != 0 {
		// We append meta info so caller can detect misuse later if desired.
		kv = append(kv, "_syserrors_kv_error", "odd number of key/value arguments")
	}
	return kv
}
