package constants

import "time"

const (
	DRIVER_NAME_POSTGRES = "postgres"
	LOGGER_NAME          = "NAME"
)

// Environment variables
const (
	APP_ENV                  = "APP_ENV"
	APP_ROOT_PATH_ENV        = "APP_ROOT_PATH_ENV"
	CONFIG_PATH_ENV          = "CONFIG_PATH_ENV"
	DB_CONN_MAX_FILETIME_ENV = "DB_CONN_MAX_FILETIME_ENV"
	DB_MAX_IDLE_CONNS_ENV    = "DB_MAX_IDLE_CONNS_ENV"
	DB_MAX_OPEN_CONNS_ENV    = "DB_MAX_OPEN_CONNS_ENV"
	LOG_CONFIG_LOG_TYPE_ENV  = "LOG_CONFIG_LOG_TYPE_ENV"
)

const (
	APP_START_TIMEOUT = 20 * time.Second
	APP_STOP_TIMEOUT  = 20 * time.Second
)

const (
	HEADER_REQUEST_ID = "X-Request-ID"
)
