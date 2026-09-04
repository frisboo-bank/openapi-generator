package contracts

import "context"

type (
	Request      any
	Command      Request
	Query        Request
	Notification any

	Handler[TRequest Request, TResponse any] interface {
		Handle(ctx context.Context, request TRequest) (TResponse, error)
	}

	NotificationHandler[TNotification Notification] interface {
		Handle(ctx context.Context, notification TNotification) error
	}

	Unit struct{}

	Mediator interface {
		Send(ctx context.Context, request any) (any, error)
		Publish(ctx context.Context, notification any) error
	}
)
