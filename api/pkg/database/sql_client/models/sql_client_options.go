package models

import (
	"fmt"
	"time"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_ssl_mode"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/models/enums/sql_client_type"
)

var _ configContracts.Configurable = (*SQLClientOptions)(nil)

type SQLClientOptions struct {
	IsEnabled             bool                              `mapstructure:"enabled"`
	Type                  sqlclienttype.SqlClientType       `mapstructure:"type"`
	Debug                 bool                              `mapstructure:"debug"`
	Host                  string                            `mapstructure:"host"`
	Port                  string                            `mapstructure:"port"`
	Database              string                            `mapstructure:"database"`
	User                  string                            `mapstructure:"user"`
	Password              string                            `mapstructure:"password"`
	SSLMode               sqlclientsslmode.SqlClientSSLMode `mapstructure:"sslMode"`
	EnableTracing         bool                              `mapstructure:"enableTracing"`
	ConnectionTimeout     time.Duration                     `mapstructure:"connectionTimeout"`
	MaxOpenConnections    int                               `mapstructure:"maxOpenConns"`
	MaxIdleConnections    int                               `mapstructure:"maxIdleConns"`
	ConnectionMaxLifetime time.Duration                     `mapstructure:"connMaxLifetime"`
	ConnectionMaxIdleTime time.Duration                     `mapstructure:"connMaxIdleTime"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (o *SQLClientOptions) GetEnabled() bool  { return o.IsEnabled }
func (o *SQLClientOptions) GetLogger() string { return o.Logger }

func (o *SQLClientOptions) SetDefaults() {
	if o.SSLMode == sqlclientsslmode.SqlClientSSLModes.UNKNOWN {
		o.SSLMode = sqlclientsslmode.SqlClientSSLModes.DISABLED
	}
	if o.ConnectionTimeout <= 0 {
		o.ConnectionTimeout = 10 * time.Second
	}
	if o.MaxOpenConnections <= 0 {
		o.MaxOpenConnections = 25
	}
	if o.MaxIdleConnections <= 0 {
		o.MaxIdleConnections = 10
	}
	if o.ConnectionMaxLifetime <= 0 {
		o.ConnectionMaxLifetime = 5 * time.Minute
	}
	if o.ConnectionMaxIdleTime <= 0 {
		o.ConnectionMaxIdleTime = 1 * time.Minute
	}
}

func (o *SQLClientOptions) Validate() error {
	if !o.IsEnabled {
		return nil
	}
	if !o.Type.IsValid() {
		return fmt.Errorf("client type is invalid")
	}
	if o.Host == "" {
		return fmt.Errorf("host is required")
	}
	if o.Port == "" {
		return fmt.Errorf("port is required")
	}
	if o.Database == "" {
		return fmt.Errorf("database is required")
	}
	if o.User == "" {
		return fmt.Errorf("user is required")
	}
	return nil
}
