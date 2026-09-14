package validation

import (
	"fmt"
	"reflect"
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
		nerr = syserrors.Newf("assert err can only be an error or a string: get %v\n", err)
	}

	panic(nerr)
}

func AssertNotNil(name string, value any) {
	if value == nil {
		panic(syserrors.CantBeNilError(name))
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		if v.IsNil() {
			panic(syserrors.CantBeNilError(name))
		}
	}
}

func AssertNotEmpty(name string, value string) {
	if strings.TrimSpace(value) == "" {
		panic(syserrors.CantBeEmptyError(name))
	}
}

func AssertValidEnum(name string, value interface {
	IsValid() bool
	String() string
},
) {
	if !value.IsValid() {
		panic(fmt.Errorf("%s is invalid: got %v", name, value.String()))
	}
}
