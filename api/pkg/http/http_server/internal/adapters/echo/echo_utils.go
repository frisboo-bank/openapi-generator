package echo

import (
	"net/http"

	"frisboo-bank/openapi-generator-service/pkg/syserrors"

	echoVendor "github.com/labstack/echo/v4"
)

func ToMiddlewaresType(middlewares ...any) ([]echoVendor.MiddlewareFunc, error) {
	var result []echoVendor.MiddlewareFunc

	for i, m := range middlewares {
		if m == nil {
			continue
		}
		mw, err := ToMiddlewareType(m)
		if err != nil {
			return nil, syserrors.Wrapf(err, "error while converting middleware %d", i)
		}

		result = append(result, mw)
	}

	return result, nil
}

func ToMiddlewareType(middleware any) (echoVendor.MiddlewareFunc, error) {
	switch mw := middleware.(type) {
	case nil:
		return func(next echoVendor.HandlerFunc) echoVendor.HandlerFunc { return next }, nil
	case echoVendor.MiddlewareFunc:
		return mw, nil
	case func(echoVendor.HandlerFunc) echoVendor.HandlerFunc:
		return echoVendor.MiddlewareFunc(mw), nil
	}
	return nil, syserrors.New(
		"invalid middleware type, must be echo.MiddlewareFunc or func(echo.HandlerFunc) echo.HandlerFunc",
	)
}

// ToHandlerFunc converts a generic handler to echo.HandlerFunc.
func ToHandlerFunc(handler any) (echoVendor.HandlerFunc, error) {
	switch h := handler.(type) {
	case echoVendor.HandlerFunc:
		return h, nil
	case func(echoVendor.Context) error:
		return h, nil
	case http.HandlerFunc:
		return echoVendor.WrapHandler(h), nil
	case http.Handler:
		return echoVendor.WrapHandler(h), nil
	}
	return nil, syserrors.New(
		"invalid handler type, must be echo.HandlerFunc, func(echo.Context) error, http.Handler, or http.HandlerFunc",
	)
}
