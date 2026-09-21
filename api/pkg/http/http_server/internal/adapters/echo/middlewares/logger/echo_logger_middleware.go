package logger

import (
	"fmt"
	"time"

	loggerContracts "frisboo-bank/openapi-generator-service/pkg/logger/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoLoggerMiddlewareOptions struct {
	Skipper middleware.Skipper
}

func NewEchoLoggerMiddleware(logger loggerContracts.Logger, config EchoLoggerMiddlewareOptions) echo.MiddlewareFunc {
	validation.AssertNotNil("logger", logger)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			stop := time.Now()

			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()
			fields := loggerContracts.Fields{
				"remote_ip":  c.RealIP(),
				"latency":    stop.Sub(start).String(),
				"host":       req.Host,
				"request":    fmt.Sprintf("%s %s", req.Method, req.RequestURI),
				"status":     res.Status,
				"size":       res.Size,
				"user_agent": req.UserAgent(),
			}
			if id := req.Header.Get(echo.HeaderXRequestID); id != "" {
				fields["request_id"] = id
			} else if id = res.Header().Get(echo.HeaderXRequestID); id != "" {
				fields["request_id"] = id
			}

			switch {
			case res.Status >= 500:
				logger.Errorw("Server error", fields)
			case res.Status >= 400:
				logger.Errorw("Client error", fields)
			case res.Status >= 300:
				logger.Infow("Redirection", fields)
			default:
				logger.Infow("Success", fields)
			}
			return nil
		}
	}
}
