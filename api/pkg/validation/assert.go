package validation

import (
	"strings"

	"frisboo-bank/openapi-generator-service/pkg/syserrors"
)

func Assert(condition bool, err any, prefix ...string) {
	if condition {
		return
	}

	var nerr error
	switch err := err.(type) {
	case error:
		nerr = err
	case string:
		nerr = syserrors.New(err)
	default:
		panic(syserrors.Newf("assert err can only be an error or a string: get %v\n", err))
	}

	panic(nerr)
}

func AssertNotNil(name string, value any) {
	Assert(value != nil, syserrors.CantBeNilError(name))
}

func AssertNotEmpty(name string, value string) {
	Assert(strings.TrimSpace(value) != "", syserrors.CantBeEmptyError(name))
}
