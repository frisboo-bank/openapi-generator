package contracts

import (
	"context"

	applicationerrorcontracts "frisboo-bank/openapi-generator-service/pkg/application_error/contracts"
)

type (
	TRequest      = any
	Request       = any
	Response      = any
	TNotification = any
	Notification  = any

	RequestHandler[TRequest Request, TResponse Response] interface {
		Handle(ctx context.Context, request TRequest) (TResponse, applicationerrorcontracts.AppError)
	}

	NotificationHandler[TNotification Notification] interface {
		Handle(ctx context.Context, notification TNotification) applicationerrorcontracts.AppError
	}

	Mediator interface {
		Send(ctx context.Context, request Request) (Response, error)
		Publish(ctx context.Context, notification Notification) error
		RegisterRequest(requestType TRequest, handler any) error
		RegisterNotification(notificationType TNotification, handler any) error
	}
)
