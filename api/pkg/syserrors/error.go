package syserrors

import (
	verr "gitlab.com/tozd/go/errors"
)

type E = verr.E

func New(message string, kv ...any) E {
	return verr.WithDetails(verr.New(message), kv...)
}

func Newf(format string, args ...any) E {
	return verr.Errorf(format, args...)
}

func Wrap(err error, message string, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(verr.Wrap(err, message), ensureDetails(kv...)...)
}

func Wrapf(err error, format string, args ...any) E {
	return verr.Wrapf(err, format, args...)
}

func Message(err error, prefix []string, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(verr.WithMessage(err, prefix...), kv...)
}

func Messagef(err error, format string, args []any, kv ...any) E {
	if err == nil {
		return nil
	}
	return verr.WithDetails(
		verr.WithMessagef(err, format, args...),
		kv...,
	)
}

func WithDetails(err error, kv ...any) E {
	return verr.WithDetails(err, kv...)
}

func Join(errs ...error) E {
	return verr.Join(errs...)
}

func Prefix(err error, prefix ...error) E {
	return verr.Prefix(err, prefix...)
}

func WithStack(err error) E {
	return verr.WithStack(err)
}

func WrapWith(err error, with error) E {
	return verr.WrapWith(err, with)
}

func Cause(err error) error {
	return verr.Cause(err)
}

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

func ensureDetails(kv ...any) []any {
	if len(kv)%2 != 0 {
		// We append meta info so caller can detect misuse later if desired.
		kv = append(kv, "_syserrors_kv_error", "odd number of key/value arguments")
	}
	return kv
}
