package config

import (
	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4/middleware"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Skipper middleware.Skipper
}

type Option = options.OptionFn[Config]

func Default() Config {
	return Config{
		Skipper: middleware.DefaultSkipper,
	}
}

func New(opts ...Option) (Config, error) {
	var cfg Config

	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return cfg, err
	}

	return base, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Skipper, validation.Required),
	)
}

var Skipper = options.Option(func(c *Config, skipper middleware.Skipper) {
	c.Skipper = skipper
})
