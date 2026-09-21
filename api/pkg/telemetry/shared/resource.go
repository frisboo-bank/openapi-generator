package shared

import (
	"context"
	"sync"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	resourceOnce sync.Once
	resource     *sdkresource.Resource
)

// Resource returns the shared process-wide OTLP resource, built once.
func Resource(serviceName string) *sdkresource.Resource {
	resourceOnce.Do(func() {
		r, err := sdkresource.New(
			context.Background(),
			sdkresource.WithAttributes(semconv.ServiceName(serviceName)),
		)
		if err != nil {
			resource = sdkresource.Default()
			return
		}
		resource = r
	})
	return resource
}
