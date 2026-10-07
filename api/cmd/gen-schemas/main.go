package main

import (
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	cacheconfig "frisboo-bank/openapi-generator-service/pkg/cache/config"
	appconfig "frisboo-bank/openapi-generator-service/pkg/config/config"
	migrationconfig "frisboo-bank/openapi-generator-service/pkg/database/migration/config"
	sqlclientconfig "frisboo-bank/openapi-generator-service/pkg/database/sql_client/config"
	httpconfig "frisboo-bank/openapi-generator-service/pkg/http/http_server/config"
	loggerconfig "frisboo-bank/openapi-generator-service/pkg/logger/config"
	mediatorconfig "frisboo-bank/openapi-generator-service/pkg/mediator/config"
	rpcconfig "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/config"
	logconfig "frisboo-bank/openapi-generator-service/pkg/telemetry/log/config"
	metricsconfig "frisboo-bank/openapi-generator-service/pkg/telemetry/metrics/config"
	tracerconfig "frisboo-bank/openapi-generator-service/pkg/telemetry/tracer/config"
	waiterconfig "frisboo-bank/openapi-generator-service/pkg/waiter/config"
)

// genEntry pairs a config package directory with the struct to reflect.
type genEntry struct {
	dir  string
	v    any
	name string
}

func main() {
	entries := []genEntry{
		{dir: "pkg/cache/config", v: &cacheconfig.CacheOptions{}, name: "CacheOptions"},
		{dir: "pkg/config/config", v: &appconfig.AppOptions{}, name: "AppOptions"},
		{dir: "pkg/database/migration/config", v: &migrationconfig.MigrationOptions{}, name: "MigrationOptions"},
		{dir: "pkg/database/sql_client/config", v: &sqlclientconfig.SQLClientOptions{}, name: "SQLClientOptions"},
		{dir: "pkg/http/http_server/config", v: &httpconfig.HTTPServerOptions{}, name: "HTTPServerOptions"},
		{dir: "pkg/logger/config", v: &loggerconfig.LoggerOptions{}, name: "LoggerOptions"},
		{dir: "pkg/mediator/config", v: &mediatorconfig.MediatorOptions{}, name: "MediatorOptions"},
		{dir: "pkg/rpc/rpc_server/config", v: &rpcconfig.RPCServerOptions{}, name: "RPCServerOptions"},
		{dir: "pkg/telemetry/log/config", v: &logconfig.LogOptions{}, name: "LogOptions"},
		{dir: "pkg/telemetry/metrics/config", v: &metricsconfig.MetricsOptions{}, name: "MetricsOptions"},
		{dir: "pkg/telemetry/tracer/config", v: &tracerconfig.TracerOptions{}, name: "TracerOptions"},
		{dir: "pkg/waiter/config", v: &waiterconfig.WaiterOptions{}, name: "WaiterOptions"},
	}

	reflector := new(jsonschema.Reflector)
	for _, e := range entries {
		schema := reflector.Reflect(e.v)
		out, err := schema.MarshalJSON()
		if err != nil {
			panic(err)
		}
		dir := filepath.Join("api", e.dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			panic(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "schema.json"), out, 0o644); err != nil {
			panic(err)
		}
	}
}
