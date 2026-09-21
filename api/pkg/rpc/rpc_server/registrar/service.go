package registrar

import "context"

type (
	Service interface {
		Name() string
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		HealthStatus(ctx context.Context) any
	}
)
