package models

import (
	"fmt"
	"net"
	"time"

	cachetype "frisboo-bank/openapi-generator-service/pkg/cache/models/enums/cache_type"
	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
)

var _ configContracts.Configurable = (*CacheOptions)(nil)

type CacheOptions struct {
	IsEnabled    bool                `mapstructure:"enabled"`
	Type         cachetype.CacheType `mapstructure:"type"`
	Host         string              `mapstructure:"host"`
	Port         string              `mapstructure:"port"`
	Password     string              `mapstructure:"password"`
	DB           int                 `mapstructure:"db"`
	PoolSize     int                 `mapstructure:"poolSize"`
	MinIdleConns int                 `mapstructure:"minIdleConns"`
	MaxRetries   int                 `mapstructure:"maxRetries"`
	DialTimeout  time.Duration       `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration       `mapstructure:"readTimeout"`
	WriteTimeout time.Duration       `mapstructure:"writeTimeout"`

	// Memory-specific
	MaxEntries int64 `mapstructure:"maxEntries"`

	// Dependencies
	Logger string `mapstructure:"logger"`
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
