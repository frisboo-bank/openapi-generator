package models

import (
	"fmt"
	"net"
	"strings"
	"time"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/models/enums/http_server_type"
)

var _ configContracts.Configurable = (*HTTPServerOptions)(nil)

type HTTPServerOptions struct {
	IsEnabled             bool                          `mapstructure:"enabled"`
	Type                  httpservertype.HttpServerType `mapstructure:"type"`
	Debug                 bool                          `mapstructure:"debug"`
	Mode                  string                        `mapstructure:"mode"`
	Host                  string                        `mapstructure:"host"`
	Port                  string                        `mapstructure:"port"`
	BasePath              string                        `mapstructure:"basePath"`
	IgnoreLogUrls         []string                      `mapstructure:"ignoreLogUrls"`
	TrustedProxies        []string                      `mapstructure:"trustedProxies"`
	MaxHeaderBytes        int                           `mapstructure:"maxHeaderBytes"`
	BodyLimit             string                        `mapstructure:"bodyLimit"`
	IdleTimeout           time.Duration                 `mapstructure:"idleTimeout"`
	ReadHeaderTimeout     time.Duration                 `mapstructure:"readHeaderTimeout"`
	ReadTimeout           time.Duration                 `mapstructure:"readTimeout"`
	ServerShutdownTimeout time.Duration                 `mapstructure:"serverShutdownTimeout"`
	WriteTimeout          time.Duration                 `mapstructure:"writeTimeout"`
	GzipLevel             int                           `mapstructure:"gzipLevel"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c *HTTPServerOptions) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *HTTPServerOptions) GetEnabled() bool  { return c.IsEnabled }
func (c *HTTPServerOptions) GetLogger() string { return c.Logger }

func (c *HTTPServerOptions) SetDefaults() {
	if c.Mode == "" {
		c.Mode = "release"
	}
	if c.BasePath == "" {
		c.BasePath = "/"
	}
	if c.BodyLimit == "" {
		c.BodyLimit = "2M"
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = 120 * time.Second
	}
	if c.ReadHeaderTimeout == 0 {
		c.ReadHeaderTimeout = 5 * time.Second
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = 30 * time.Second
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = 30 * time.Second
	}
	if c.ServerShutdownTimeout == 0 {
		c.ServerShutdownTimeout = 30 * time.Second
	}
	if c.GzipLevel == 0 {
		c.GzipLevel = 5
	}
	if c.MaxHeaderBytes == 0 {
		c.MaxHeaderBytes = 8 * 1024
	}
}

func (c *HTTPServerOptions) Validate() error {
	if !c.IsEnabled {
		return nil
	}
	if !c.Type.IsValid() {
		return fmt.Errorf("invalid server type")
	}
	if strings.TrimSpace(c.Host) == "" {
		return fmt.Errorf("host is required")
	}
	if strings.TrimSpace(c.Port) == "" {
		return fmt.Errorf("port is required")
	}
	return nil
}
