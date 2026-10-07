package config

import (
	"fmt"
	"time"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	sqlclientsslmode "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclientsslmode"
	sqlclienttype "frisboo-bank/openapi-generator-service/pkg/database/sql_client/types/sqlclienttype"
)

var _ configContracts.Configurable = (*SQLClientOptions)(nil)

type SQLClientOptions struct {
	IsEnabled bool `mapstructure:"enabled" json:"enabled"`
	Type sqlclienttype.SqlClientType `mapstructure:"type" json:"type"`
	Debug bool `mapstructure:"debug" json:"debug"`
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
	Database string `mapstructure:"database" json:"database"`
	User string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	SSLMode sqlclientsslmode.SqlClientSSLMode `mapstructure:"sslMode" json:"sslMode"`
	EnableTracing bool `mapstructure:"enableTracing" json:"enableTracing"`
	EnableMetrics bool `mapstructure:"enableMetrics" json:"enableMetrics"`
	ConnectionTimeout time.Duration `mapstructure:"connectionTimeout" json:"connectionTimeout"`
	MaxOpenConnections int `mapstructure:"maxOpenConns" json:"maxOpenConns"`
	MaxIdleConnections int `mapstructure:"maxIdleConns" json:"maxIdleConns"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime"`
	ConnectionMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime" json:"connMaxIdleTime"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
	Tracer string `mapstructure:"tracer" json:"tracer"`
	Metrics string `mapstructure:"metrics" json:"metrics"`
}

func (o *SQLClientOptions) Enable() bool { return o.IsEnabled }
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
	// sqlite3x is file-based: Database is the file path; no host/port/user.
	if o.Type == sqlclienttype.SqlClientTypes.SQLITE3X {
		if o.Database == "" {
			return fmt.Errorf("database (file path) is required")
		}
		return nil
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
