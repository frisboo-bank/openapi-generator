package mediator_test

import (
	"context"
	"frisboo-bank/openapi-generator-service/pkg/logger"
	"frisboo-bank/openapi-generator-service/pkg/mediator/config"
	"testing"

	environmentenum "frisboo-bank/openapi-generator-service/pkg/environment/models/enums/environment"
	"frisboo-bank/openapi-generator-service/pkg/mediator"
	"frisboo-bank/openapi-generator-service/pkg/mediator/contracts"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRequest struct {
	Value string
}

type testResponse struct {
	Value string
}

type testHandler struct {
	onHandle func(context.Context, *testRequest) (*testResponse, error)
}

func handler() testHandler {
	return testHandler{onHandle: func(_ context.Context, req *testRequest) (*testResponse, error) {
		return &testResponse{Value: req.Value}, nil
	}}
}

func (h testHandler) Handle(ctx context.Context, req *testRequest) (*testResponse, error) {
	return h.onHandle(ctx, req)
}

type testNotification struct {
	Value string
}

func createNewMediator(t *testing.T) contracts.Mediator {
	t.Helper()

	mediator, err := mediator.NewMediator(
		&config.MediatorOptions{},
		logger.CreateNoopLogger("test", environmentenum.Environments.DEVELOPMENT),
	)
	require.NoError(t, err)

	return mediator
}

func TestSendRequest(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	require.NoError(t, mediator.RegisterRequest(mediatorInstance, handler()))

	res, err := mediator.Send[testRequest, testResponse](mediatorInstance, context.Background(), &testRequest{
		Value: "test",
	})

	require.NoError(t, err)
	assert.Equal(t, "test", res.Value)
}

func TestSendRequestError(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	handler := testHandler{
		onHandle: func(_ context.Context, req *testRequest) (*testResponse, error) {
			return nil, assert.AnError
		},
	}

	err := mediator.RegisterRequest(mediatorInstance, handler)
	require.NoError(t, err)

	_, err = mediator.Send[testRequest, testResponse](mediatorInstance, context.Background(), &testRequest{
		Value: "test",
	})

	assert.Error(t, err)
}

func TestSendNoHandlerRequest(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	_, err := mediator.Send[testRequest, testResponse](mediatorInstance, context.Background(), &testRequest{
		Value: "test",
	})

	assert.Error(t, err)
}

func TestSendCancelledContextReturnsError(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	require.NoError(t, mediator.RegisterRequest(mediatorInstance, handler()))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := mediator.Send[testRequest, testResponse](mediatorInstance, ctx, &testRequest{
		Value: "test",
	})

	assert.Error(t, err)
}

func TestRegisterDuplicateRequest(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	h := handler()

	require.NoError(t, mediator.RegisterRequest(mediatorInstance, h))
	assert.Error(t, mediator.RegisterRequest(mediatorInstance, h))
}

func TestSendPointerResolvesValueKey(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	require.NoError(t, mediatorInstance.RegisterRequest(testRequest{}, func(_ context.Context, req any) (any, error) {
		r, ok := req.(*testRequest)
		require.True(t, ok)
		return &testResponse{Value: r.Value}, nil
	}))

	res, err := mediatorInstance.Send(context.Background(), &testRequest{Value: "ptr"})

	require.NoError(t, err)
	assert.Equal(t, "ptr", res.(*testResponse).Value)
}

func TestSendNilRequestPanics(t *testing.T) {
	mediatorInstance := createNewMediator(t)

	assert.Panics(t, func() {
		_, _ = mediatorInstance.Send(context.Background(), nil)
	})
}
