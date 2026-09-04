package models

import (
	"fmt"
	"net"
	"strconv"
	"time"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	rpcservertype "frisboo-bank/openapi-generator-service/pkg/rpc/rpc_server/models/enums/rpc_server_type"
)

var _ configContracts.Configurable = (*RPCServerOptions)(nil)

type RPCServerOptions struct {
	IsEnabled             bool                        `mapstructure:"enabled"`
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
	Debug                 bool                        `mapstructure:"debug"`
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	KeepAliveTime         time.Duration               `mapstructure:"keepAliveTime"`
	KeepAliveTimeout      time.Duration               `mapstructure:"KeepAliveTimeout"`
	MaxConnectionAge      time.Duration               `mapstructure:"MaxConnectionAge"`
	MaxConnectionAgeGrace time.Duration               `mapstructure:"MaxConnectionAgeGrace"`
	MaxConnectionIdle     time.Duration               `mapstructure:"MaxConnectionIdle"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c *RPCServerOptions) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *RPCServerOptions) GetEnabled() bool  { return c.IsEnabled }
func (c *RPCServerOptions) GetLogger() string { return c.Logger }

func (c *RPCServerOptions) SetDefaults() {
	if c.KeepAliveTime == 0 {
		c.KeepAliveTime = 10 * time.Minute
	}
	if c.KeepAliveTimeout == 0 {
		c.KeepAliveTimeout = 20 * time.Second
	}
	if c.MaxConnectionIdle == 0 {
		c.MaxConnectionIdle = 5 * time.Minute
	}
	if c.MaxConnectionAge == 0 {
		c.MaxConnectionAge = 5 * time.Minute
	}
	if c.MaxConnectionAgeGrace == 0 {
		c.MaxConnectionAgeGrace = 1 * time.Minute
	}
}

func (c *RPCServerOptions) Validate() error {
	if !c.IsEnabled {
		return nil
	}
	if !c.Type.IsValid() {
		return fmt.Errorf("invalid server type")
	}
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	port, err := strconv.Atoi(c.Port)
	if err != nil || (port <= 0 || port > 65535) {
		return fmt.Errorf("port must be a valid port number")
	}
	if c.KeepAliveTimeout > c.KeepAliveTime && c.KeepAliveTime > 0 {
		return fmt.Errorf("keepAliveTimeout (%v) cannot exceed keepAliveTime (%v)", c.KeepAliveTimeout, c.KeepAliveTime)
	}
	return nil
}
