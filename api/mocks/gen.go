package mocks

//go:generate mockgen -source ../pkg/application/contracts/application.go -destination pkg/application/contracts/application.go -package mocks
//go:generate mockgen -source ../pkg/application/contracts/application_builder.go -destination pkg/application/contracts/application_builder.go -package mocks
//go:generate mockgen -source ../pkg/cache/contracts/cache.go -destination pkg/cache/contracts/cache.go -package mocks
//go:generate mockgen -source ../pkg/config/contracts/config_loader.go -destination pkg/config/contracts/config_loader.go -package mocks
//go:generate mockgen -source ../pkg/database/sql_client/contracts/sql_client.go -destination pkg/database/sql_client/contracts/sql_client.go -package mocks
//go:generate mockgen -source ../pkg/http/http_server/contracts/http_server.go -destination pkg/http/http_server/contracts/http_server.go -package mocks
//go:generate mockgen -source ../pkg/logger/contracts/logger.go -destination pkg/logger/logger_mock.go -package mocks
//go:generate mockgen -source ../pkg/mediator/contracts/mediator.go -destination pkg/mediator/contracts/mediator.go -package mocks
//go:generate mockgen -source ../pkg/rpc/rpc_server/contracts/rpc_server.go -destination pkg/rpc/rpc_server/contracts/rpc_server.go -package mocks
