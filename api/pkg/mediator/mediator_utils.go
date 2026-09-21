package mediator

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"
	"frisboo-bank/openapi-generator-service/pkg/validation"
)

func Send[TRequest contracts.Request, TResponse contracts.Response](
	mediator contracts.Mediator,
	ctx context.Context,
	request *TRequest,
) (*TResponse, error) {
	validation.AssertNotNil("mediator", mediator)
	validation.AssertNotNil("request", request)

	res, err := mediator.Send(ctx, request)
	if err != nil {
		return nil, err
	}

	typedRes, ok := res.(*TResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: got %T, want %T", res, new(TResponse))
	}

	return typedRes, nil
}

func Publish[TNotification contracts.Notification](
	mediator contracts.Mediator,
	ctx context.Context,
	notification *TNotification,
) error {
	validation.AssertNotNil("mediator", mediator)
	validation.AssertNotNil("notification", notification)

	return mediator.Publish(ctx, notification)
}

func RegisterRequest[TRequest contracts.Request, TResponse contracts.Response](
	mediator contracts.Mediator,
	handler contracts.RequestHandler[TRequest, TResponse],
) error {
	validation.AssertNotNil("mediator", mediator)
	validation.AssertNotNil("handler", handler)

	wrapper := func(ctx context.Context, request contracts.Request) (contracts.Response, error) {
		typedReq, ok := request.(TRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type: got %T, want %T", request, new(TRequest))
		}
		return handler.Handle(ctx, typedReq)
	}

	var zeroReq TRequest
	return mediator.RegisterRequest(zeroReq, wrapper)
}

func RegisterNotification[TNotification contracts.Notification](
	mediator contracts.Mediator,
	handler contracts.NotificationHandler[TNotification],
) error {
	validation.AssertNotNil("mediator", mediator)
	validation.AssertNotNil("handler", handler)

	wrapper := func(ctx context.Context, request contracts.Request) error {
		typedReq, ok := request.(TNotification)
		if !ok {
			return fmt.Errorf("unexpected request type: got %T, want %T", request, new(TNotification))
		}
		return handler.Handle(ctx, typedReq)
	}

	var zeroReq TNotification
	return mediator.RegisterNotification(zeroReq, wrapper)
}
