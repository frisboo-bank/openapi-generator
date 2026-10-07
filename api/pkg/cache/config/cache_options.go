package config

import (
	"fmt"
	"net"
	"time"

	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/types/cachetype"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
)

var _ configContracts.Configurable = (*CacheOptions)(nil)

type CacheOptions struct {
	IsEnabled bool `mapstructure:"enabled" json:"enabled"`
	Type cachetype.CacheType `mapstructure:"type" json:"type"`
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	DB int `mapstructure:"db" json:"db"`
	PoolSize int `mapstructure:"poolSize" json:"poolSize"`
	MinIdleConns int `mapstructure:"minIdleConns" json:"minIdleConns"`
	MaxRetries int `mapstructure:"maxRetries" json:"maxRetries"`
	DialTimeout time.Duration `mapstructure:"dialTimeout" json:"dialTimeout"`
	ReadTimeout time.Duration `mapstructure:"readTimeout" json:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`

	// Memory-specific
	MaxEntries int64 `mapstructure:"maxEntries" json:"maxEntries"`

	// Dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

func (c *CacheOptions) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *CacheOptions) GetEnabled() bool  { return c.IsEnabled }
func (c *CacheOptions) GetLogger() string { return c.Logger }

func (c *CacheOptions) SetDefaults() {
	if c.PoolSize == 0 {
		c.PoolSize = 10
	}
	if c.MinIdleConns == 0 {
		c.MinIdleConns = 5
	}
	if c.MaxRetries == 0 {
		c.MaxRetries = 3
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = 5 * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 3 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 3 * time.Second
	}
	if c.MaxEntries == 0 {
		c.MaxEntries = 10_000
	}
}

func (c *CacheOptions) Validate() error {
	if !c.IsEnabled {
		return nil
	}
	if !c.Type.IsValid() {
		return fmt.Errorf("invalid cache type")
	}
	if c.Host == "" && c.Type != cachetype.CacheTypes.MEMORY {
		return fmt.Errorf("host is required")
	}
	if c.Port == "" && c.Type != cachetype.CacheTypes.MEMORY {
		return fmt.Errorf("port is required")
	}
	return nil
}

//go:generate go run github.com/invopop/jsonschema -o schema.json -package config CacheOptions
