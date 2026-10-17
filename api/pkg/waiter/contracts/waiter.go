package contracts

import (
	"context"
	"errors"

	"frisboo-bank/openapi-generator-service/pkg/waiter/types/hookstatus"
)

var ErrForcedStop = errors.New("waiter: forced stop")

type (
	WaiterHook struct {
		Name   string
		Status func() hookstatus.HookStatus
		Start  func(ctx context.Context) error
		Stop   func(ctx context.Context) error
	}

	Waiter interface {
		AddHooks(hooks ...WaiterHook) error
		AddHook(hook WaiterHook) error
		Status() map[string]hookstatus.HookStatus
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}
)

func (h WaiterHook) IsRunning() bool {
	return h.Status == nil || h.Status() == hookstatus.HookStatuses.RUNNING
}
