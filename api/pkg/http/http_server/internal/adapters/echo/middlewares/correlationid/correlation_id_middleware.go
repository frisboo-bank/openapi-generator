package correlationid

import (
	"context"

	loggercontracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ctxKey string

const correlationIDKey = ctxKey("correlationID")

func NewEchoLoggerMiddleware(logger loggercontracts.Logger) echo.MiddlewareFunc {
	validation.AssertNotNil("logger", logger)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			id := req.Header.Get(echo.HeaderXCorrelationID)
			if id == "" {
				id = uuid.New().String()
			}

			c.Response().Header().Set(echo.HeaderXCorrelationID, id)

			newCtx := context.WithValue(req.Context(), correlationIDKey, id)
			c.SetRequest(req.WithContext(newCtx))

			return next(c)
		}
	}
}
