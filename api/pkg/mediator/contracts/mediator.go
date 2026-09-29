package contracts

import "context"

type (
	TRequest      = any
	Request       = any
	Response      = any
	TNotification = any
	Notification  = any

	RequestHandler[TRequest Request, TResponse Response] interface {
		Handle(ctx context.Context, request TRequest) (TResponse, error)
	}

	NotificationHandler[TNotification Notification] interface {
		Handle(ctx context.Context, notification TNotification) error
	}

	Mediator interface {
		Send(ctx context.Context, request Request) (Response, error)
		Publish(ctx context.Context, notification Notification) error
		RegisterRequest(requestType TRequest, handler any) error
		RegisterNotification(notificationType TNotification, handler any) error
	}
)
