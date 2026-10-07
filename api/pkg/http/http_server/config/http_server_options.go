package config

import (
	"fmt"
	"net"
	"strings"
	"time"

	configContracts "frisboo-bank/openapi-generator-service/pkg/config/contracts"
	httpservertype "frisboo-bank/openapi-generator-service/pkg/http/http_server/types/httpservertype"
)

var _ configContracts.Configurable = (*HTTPServerOptions)(nil)

type HTTPServerOptions struct {
	IsEnabled bool `mapstructure:"enabled" json:"enabled"`
	Type httpservertype.HttpServerType `mapstructure:"type" json:"type"`
	Debug bool `mapstructure:"debug" json:"debug"`
	Mode string `mapstructure:"mode" json:"mode"`
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
	BasePath string `mapstructure:"basePath" json:"basePath"`
	IgnoreLogUrls []string `mapstructure:"ignoreLogUrls" json:"ignoreLogUrls"`
	TrustedProxies []string `mapstructure:"trustedProxies" json:"trustedProxies"`
	MaxHeaderBytes int `mapstructure:"maxHeaderBytes" json:"maxHeaderBytes"`
	BodyLimit string `mapstructure:"bodyLimit" json:"bodyLimit"`
	IdleTimeout time.Duration `mapstructure:"idleTimeout" json:"idleTimeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout" json:"readHeaderTimeout"`
	ReadTimeout time.Duration `mapstructure:"readTimeout" json:"readTimeout"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout" json:"serverShutdownTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout"`
	GzipLevel int `mapstructure:"gzipLevel" json:"gzipLevel"`

	// dependencies
	Logger string `mapstructure:"logger" json:"logger"`
}

func (c *HTTPServerOptions) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func (c *HTTPServerOptions) Enable() bool  { return c.IsEnabled }
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
