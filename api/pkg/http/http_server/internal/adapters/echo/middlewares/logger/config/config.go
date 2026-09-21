package config

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ validation.Validatable = (*Config)(nil)

type Config struct {
	Skipper middleware.Skipper
}

type Option func(*Config)

func Default() Config {
	return Config{
		Skipper: middleware.DefaultSkipper,
	}
}

func New(opts ...Option) (Config, error) {
	cfg := Default()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Skipper, validation.Required),
	)
}

func Skipper(skipper middleware.Skipper) Option {
	return func(c *Config) {
		c.Skipper = skipper
	}
}
