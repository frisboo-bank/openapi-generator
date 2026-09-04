package postgres

import (
	"context"
	"fmt"

	"frisboo-bank/openapi-generator-service/pkg/database/sql_client/models"
	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_ssl_mode"

	"github.com/jmoiron/sqlx"
)

var postgresSSLModeMap = map[sqlclientsslmode.SqlClientSSLMode]string{
	sqlclientsslmode.SqlClientSSLModes.DISABLED:   "disable",
	sqlclientsslmode.SqlClientSSLModes.REQUIRE:    "require",
	sqlclientsslmode.SqlClientSSLModes.VERIFYCA:   "verify-ca",
	sqlclientsslmode.SqlClientSSLModes.VERIFYFULL: "verify-full",
}

func ConnectToPostgres(cfg *models.SQLClientOptions) (*sqlx.DB, error) {
	sslMode := postgresSSLModeMap[cfg.SSLMode]
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, sslMode,
	)

	if cfg.ConnectionTimeout > 0 {
		dsn += fmt.Sprintf(" connect_timeout=%d", int(cfg.ConnectionTimeout.Seconds()))
	}

	connectCtx := context.Background()
	if cfg.ConnectionTimeout > 0 {
		var cancel context.CancelFunc
		connectCtx, cancel = context.WithTimeout(connectCtx, cfg.ConnectionTimeout)
		defer cancel()
	}

	db, err := sqlx.ConnectContext(connectCtx, "pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}

	if cfg.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}

	if cfg.ConnectionMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	}

	if cfg.ConnectionMaxIdleTime > 0 {
		db.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTime)
	}

	return db, nil
}
